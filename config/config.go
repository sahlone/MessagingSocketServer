package config


import (
	"os"
	"encoding/json"

	"github.com/sahilahmadlone/MessagingSocketServer/logger"
)

//ServerConfig struct holds basic configurations
//These can be set by passing arguments in to the commandline
//in the form of `logLevel=DEBUG`
//All configurations passed by commandline must be space delimited
//If no commandline arguments are passed, default configs are used
type ServerConfig struct {
	LogLevel	string
	EventListenerPort	int
	ClientListenerPort	int
	SequenceNumber	int
}

//Loads default configuration for the Server from conf.json
//In a production environment such configuration would likely be done
//using feature flags or a puppet-like tool to avoid code changes upon config update.
func ServerDefaultConfig(path string) *ServerConfig {
	var conf ServerConfig = ServerConfig{}
	file, err := os.Open(path+"conf.json")
	if err != nil {
		logger.Error("FileError", err)
	}
	decoder := json.NewDecoder(file)
	decErr := decoder.Decode(&conf)
	if decErr != nil {
		logger.Error("Decode Error ", err)
	}
	logger.SetLevel(conf.LogLevel)
	return &conf
}

