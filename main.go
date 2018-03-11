package main

import (
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/sahilahmadlone/MessagingSocketServer/config"
	"github.com/sahilahmadlone/MessagingSocketServer/logger"
	"github.com/sahilahmadlone/MessagingSocketServer/server"
)

//Checks for valid parameters via commandline arguments
//Returns variable and value if appropriate
func parseEnvVar(val string) (string, string) {
	splitVal := strings.Split(val, "=")
	if len(splitVal) != 2 {
		return "", ""
	}
	envVar := splitVal[0]
	envVal := splitVal[1]
	return envVar, envVal

}

//Simple check error function
func checkError(err error) bool {
	if err != nil {
		logger.Error("Not a valid environment value ", err, " reverting to default config")
		return false
	}
	return true
}

//Updates and overrides default config values
func exportVals(env string, val string, conf config.ServerConfig) config.ServerConfig {
	switch env {
	case "logLevel":
		if ok := checkError(logger.SetLevel(val)); ok {
			conf.LogLevel = val
		}
	case "eventListenerPort":
		val, err := strconv.Atoi(val)
		if ok := checkError(err); ok {
			conf.EventListenerPort = val
		}
	case "clientListenerPort":
		val, err := strconv.Atoi(val)
		if ok := checkError(err); ok {
			conf.EventListenerPort = val
		}
	case "sequenceNumber":
		val, err := strconv.Atoi(val)
		if ok := checkError(err); ok {
			conf.EventListenerPort = val
		}
	}

	return conf

}

//Updates config according to valid commandline args
//If parsed arguments are valid
func overRideDefaultConfig(args []string, conf config.ServerConfig) *config.ServerConfig {
	for i := 1; i < len(args); i++ {
		env, val := parseEnvVar(args[1])
		if env == "" && val == "" {
			logger.Error("Invalid commandline arguments ", env, " ", val, "reverting to default config")
		} else {
			conf = exportVals(env, val, conf)

		}
	}
	return &conf
}

//Sets up configuration for Environment
//Starts the Server
//Checks for server configs by commandline
//Catches unexpected signals
func main() {
	var conf *config.ServerConfig
	conf = config.ServerDefaultConfig("config/")
	if len(os.Args) > 1 {
		conf = overRideDefaultConfig(os.Args, *conf)
	}

	logger.Info("Starting Server:")
	server, err := server.Run(*conf)
	if err != nil {
		recover()
		os.Exit(1)
	}

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	for sig := range sigChannel {
		if sig == os.Interrupt {
			logger.Info("Graceful ShutDown of Server")
			os.Exit(1)

		}
	}
	if err := server.ShutDown(); err != nil {
		logger.Error("Error shutting down Server, exiting now.")
		recover()
		os.Exit(1)
	}
}
