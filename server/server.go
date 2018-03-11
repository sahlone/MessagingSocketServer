package server

import (
	"bufio"
	"net"
	"io"
	"strconv"
	"strings"
	"errors"

	"github.com/sahilahmadlone/MessagingSocketServer/config"
	"github.com/sahilahmadlone/MessagingSocketServer/logger"
)
//Sequence counter for dispatcher
var SequenceNum int

//Server attributes for shutdown
type Server struct {
	finished      chan struct{}
	IsRunning	bool
	UListener 	net.Listener
	EListener       net.Listener
}

//Event struct for parsing and processing
type Event struct {
	sequence int
	eventType       string
	fromUserId     int
	toUserId       int
	payload string
}

//User client struct for parsing and notifying
type UserClient struct {
	userId int
	connection net.Conn
}



//Sets up the dispatcher with channels for when events start arriving
//Starts Server listening on specified ports from configuration (param)
//Starts two goroutines accepting and serving events and userClients
func Run(config config.ServerConfig) (*Server, error) {
	finished := make(chan struct{})

	SequenceNum = config.SequenceNumber

	userChannel, eventChannel, err := dispatcher(finished)

	if err != nil {
		return nil, err
	}
	es, err := net.Listen("tcp", ":"+strconv.Itoa(config.EventListenerPort))
	if err != nil {
		recover()
		return nil, err
	}
	us, err := net.Listen("tcp", ":"+strconv.Itoa(config.ClientListenerPort))
	if err != nil {
		recover()
		return nil, err
	}
	logger.Info("Listening on Ports ", strconv.Itoa(config.EventListenerPort)," and ", strconv.Itoa(config.ClientListenerPort))

	go acceptAndServeUsers(userChannel, us, finished)
	go acceptAndServeEvents(eventChannel, es, finished)
	return &Server{finished, true, us, es}, nil
}




//When listener receives event, this method handles it
//in a goroutine -- reading in the message, parsing the message, assigning values to
//Event struct, and sending `Event` to event channel
func handleEventConns(connection net.Conn, eventChan chan<- Event) {
	b := bufio.NewReader(connection)
	for {
		m, err := b.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("End of message stream", err)
				return
			}
			logger.Error(err)
			return
		}
		msg := string(m)
		msg = strings.Trim(msg, "\n")
		msg = strings.Trim(msg, "\r")
		parsedEvent, err := parseEventMessage(msg)
		if err == nil {
			eventChan <- *parsedEvent

		}
	}

}


//Similar to handling event messages, this method
//reads the message from userClient, parses the clientID,
//creates appropriate UserClient struct and sends
//`UserClient` to the user channel
func handleUserConns(connection net.Conn, userChan chan<- UserClient) {
		b := bufio.NewReader(connection)
		m, err := b.ReadString('\n')
		msg := string(m)
		msg = strings.Trim(msg, "\n")
		msg = strings.Trim(msg, "\r")
		userID, err := strconv.Atoi(msg)
		if err != nil {
			logger.Error("Bad User Request ", err)
		}

		userClient := UserClient{
			userId: userID,
			connection: connection,
		}
		userChan <- userClient

}
//Locally creates maps to keep track of received events,
// notifying userClients in the appropriate order,
// and avoid race conditions.
//Using a map implementation of a Queue in order to dispatch and processes events
// in the correct order and notifying
//all appropriate users (if connected) determined by event type
func dispatcher(finished chan struct{}) (chan<- UserClient, chan<- Event, error) {
	//Queue implementation for dispatch order
	MessageQueue := make(map[int]Event)
	//Map to keep track of followers for a given user
	FollowerMap := make(map[int]map[int]bool)
	//Map to keep track of events to users
	UserEventChannels := make(map[int]chan Event)
	//Event channel to hold events
	EChannel := make(chan Event)
	//User channel to hold clients
	UChannel := make(chan UserClient)

	go func() {
		for {
			select {
			//For incomming events
			case event := <-EChannel:
				MessageQueue[event.sequence] = event
				for {
					if event, ok := MessageQueue[SequenceNum]; ok {
						logger.Debug("SequenceNumber at ", SequenceNum, " dispatching event ", event.payload)
						delete(MessageQueue, event.sequence)
						processEventMessage(event, FollowerMap, UserEventChannels)
						SequenceNum++

					} else {
						break
					}
				}
			//For listening users
			case conUser := <-UChannel:
				evChan := make(chan Event, 1)

				var event Event
				go func() {
					for {
						select {
						case event = <-evChan:
							logger.Debug("Writing to user ", conUser.userId)
							_, err := conUser.connection.Write([]byte(event.payload+"\r"+"\n"))
							if err!= nil {
								logger.Error(err)
							}
						case <-finished:
							return
						}

					}
				}()
				UserEventChannels[conUser.userId] = evChan

			case <-finished:
				return
			}
		}

	}()
	return UChannel, EChannel, nil
}

//Takes in listener and user channel (and finished) as params
//Accepts connection and sends it to the connection channel
//In that event calls goroutine to handle user connections appropriately
func acceptAndServeUsers(userChan chan<-UserClient, listener net.Listener, finished chan struct{}) {
	for {
		connectionChannel := make(chan net.Conn)
		go func() {
			connect, err := listener.Accept()
			if err == nil {
				connectionChannel <- connect
			}
		}()

		select {
		case conChan := <-connectionChannel:
			go handleUserConns(conChan, userChan)
		case <-finished:
			listener.Close()
			return
		}
	}
}

//Similar to acceptAndServeUsers, once a connection is made it's sent to connectionChannel
//In that event the goroutine to handle and process events is started
func acceptAndServeEvents(eventChan chan<-Event, listener net.Listener, finished chan struct{}) {
	for {
		connectionChannel := make(chan net.Conn)
		go func() {
			conn, err := listener.Accept()
			if err == nil {
				connectionChannel <- conn

			}
		}()

		select {
		case conChan := <-connectionChannel:
			go handleEventConns(conChan, eventChan)
		case <-finished:
			listener.Close()
			return
		}
	}
}

//This accepts the Event itself, followers map, and a map of event conections as params
//processEventMessage logic sends the event to the appropriate channel based on the
//eventType. It will also handle all the follow/unfollow logic when needed
func processEventMessage(event Event, fm map[int]map[int]bool, eventConns map[int]chan Event) {
	logger.Debug("Processing Event ", event.payload)
	switch event.eventType {
	case "F":
		follower, ok := fm[event.toUserId]
		if !ok {
			follower = make(map[int]bool)
		}
		follower[event.fromUserId] = true
		fm[event.toUserId] = follower
		if ec, ok := eventConns[event.toUserId]; ok {
			ec <- event
		}
	case "U":
		follower := fm[event.toUserId]
		delete(follower, event.fromUserId)
	case "B":
		for _, ec := range eventConns {
			ec <- event
		}
	case "P":
		if ec, ok := eventConns[event.toUserId]; ok {
			ec <- event
		}
	case "S":
		followers := fm[event.fromUserId]
		for f, _ := range followers {
			if ec, ok := eventConns[f]; ok {
				ec <- event
			}
		}

	}
}

//Manipulates the event message received by the event listener
//First sets sequence and payload params of the `Event`
//then sets all other fields based on the type of the event
//A helper function parseUserIds is used to minimize duplicative code
//Will error, log and continue to listen with bad input.
func parseEventMessage(msg string) (*Event, error) {
	var event Event
	var err error
	event.payload = msg
	ef := strings.Split(msg, "|")
	if len(ef) < 2 || len(ef) > 4 {
		logger.Error("Bad Request", event.payload)
		return nil, errors.New("Invalid Event")
	}
	event.sequence, err = strconv.Atoi(ef[0])
	if err != nil {
		logger.Error("Bad Request", event.payload)
		return nil, err
	}
	event.eventType = ef[1]
	switch event.eventType {
	case "F":
		event.eventType = "F"
		parsedEvent, err := parseUserIds(ef, event)
		if err != nil{
			logger.Error(err)
			return nil, err
		}
		event = *parsedEvent
		return &event, nil
	case "U":
		event.eventType = "U"
		parsedEvent, err := parseUserIds(ef, event)
		if err != nil{
			logger.Error(err)
			return nil, err
		}
		event = *parsedEvent
		return &event, nil
	case "B":
		event.eventType = "B"
		return &event, nil

	case "P":
		event.eventType = "P"
		parsedEvent, err := parseUserIds(ef, event)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		event = *parsedEvent
		return &event, nil

	case "S":
		event.eventType = "S"
		fromUID, err := strconv.Atoi(ef[2])
		if err != nil {
			logger.Error("Bad Request ", err, " ", event.payload)
			return nil, err
		}
		event.fromUserId = fromUID
		return &event, nil
	}

	return &event, nil
}

//Helper function for parseEventMessage that processes the userIds
//of a given message
//Will error and continue to listen if input is bad
func parseUserIds(splitMsg []string, event Event) (*Event, error){
	if len(splitMsg) != 4 {
		logger.Error("Bad Request ", event.payload)
		err := errors.New("Invalid Request")
		return nil, err
	}
	toUID, err := strconv.Atoi(splitMsg[3])
	if err != nil {
		logger.Error("Bad Request", err, " ", event.payload)
		return nil, err
	} else {
		event.toUserId = toUID
	}
	fromUID, err := strconv.Atoi(splitMsg[2])
	if err != nil {
		logger.Error("Bad Request ", err, " ", event.payload)
		return nil, err
	}
	event.fromUserId = fromUID
	return &event, nil

}

func (ms *Server) ShutDown() error {
	close(ms.finished)
	ms.EListener.Close()
	ms.UListener.Close()
	ms.IsRunning = false
	return nil
}
