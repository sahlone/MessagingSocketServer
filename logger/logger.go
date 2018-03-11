package logger

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)
//Since only errors are logged in std golang lib this is a custom logger
//LoggingLevel can be altered in the conf.json or via commandline

var (
	lls = map[string]int{
		"ALL": 0,
		"DEBUG": 1,
		"INFO": 2,
		"WARN": 3,
		"ERROR": 4,
	}
	logLevel int
	timeFormat string
	mutex *sync.Mutex
)

//logMessage struct is created/printed during runtime when
//one of the LogLeve() functions is called
//Here the loglevel, file and line of caller (skips 3 levels), time, and message are called
type logMessage struct {
	m	[]interface{}
	file 	string
	lline	int
	level 	string
}


//Mutex is used for logging to stdOut
//Prints log messages with all the info
func writeLogMessage(lm *logMessage){
	mutex.Lock()
	timeNow := time.Now().Format(timeFormat)
	fmt.Println("[",lm.level,"] - ", timeNow, "--", lm.file,"@",lm.lline," : ", lm.m)
	mutex.Unlock()
}

//Sets the LogLevel of the logger
func SetLevel(level string) error {
	level = strings.ToUpper(level)
	if nl, ok :=lls[level]; ok {
		logLevel = nl
		return nil
	}
	return errors.New("INVALID LOG LEVEL: " + level)
}

//Sets the TimeFormat the logger should print out
//Not currently used in this program
func SetTimeFormat(tf string){
	timeFormat = tf
}

//Creates logMessage object
func makeLogMessage(m []interface{}, level string) *logMessage{
	_, file, lline, _ := runtime.Caller(3)
	return &logMessage{m, filepath.Base(file), lline, level}

}

//NOTE: If LogLevel is set to ALL, all levels of logging will be visible

//Writes the logMessage to stdOut iff logLevel<=int value of that "DEBUG"
func Debug(m ...interface{}) {
	if logLevel <= lls["DEBUG"] {
		lm := makeLogMessage(m , "DEBUG")
		writeLogMessage(lm)
	}
}
//Writes the logMessage to stdOut iff logLevel<=int value of that "ALL"
func All(m ...interface{}) {
	if logLevel <= lls["ALL"] {
		lm := makeLogMessage(m , "ALL")
		writeLogMessage(lm)
	}
}
//Writes the logMessage to stdOut iff logLevel<=int value of that "ERROR"
func Error(m ...interface{}) {
	if logLevel <= lls["ERROR"] {
		lm := makeLogMessage(m , "ERROR")
		writeLogMessage(lm)
	}
}
//Writes the logMessage to stdOut iff logLevel<=int value of that "INFO"
func Info(m ...interface{}) {
	if logLevel <= lls["INFO"] {
		lm := makeLogMessage(m , "INFO")
		writeLogMessage(lm)
	}
}
//Writes the logMessage to stdOut iff logLevel<=int value of that "WARN"
func Warn(m ...interface{}) {
	if logLevel <= lls["WARN"] {
		lm := makeLogMessage(m , "WARN")
		writeLogMessage(lm)
	}
}

func init(){
	timeFormat = "2006/01/02 - 15:00:05"
	mutex = &sync.Mutex{}
}

