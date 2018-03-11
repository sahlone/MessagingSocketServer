# Socket Server

## Dependencies
1. JDK 7 VM
2. Golang 1.8 or later

## Setup
1. Clone the repo.
2. Setup go workspace appropriately (if not already setup)
3. `cd $GOPATH/github.com/sahilahmadlone/MessagingSocketServer`
4. `go build`
5. `./MessagingSocketServer`
6. Provided tests can be used the validate the build.
6. To run tests run `./followermaze.sh`

## Server Configurations
The following parameters are configurable:
   - **logLevel**: Set log level for custom logger.
   - **clientListenerPort**: The Port the server will listen for user clients on.
   - **eventListenerPort**: The Port the server will listen for events on.
   - **sequenceNumber**: The sequence number of the first event the server should expect to receive.

These coniguration parameters can be set in the `conf.json` file, or passed in as commandline arguments to the server.
*Example:* <br />
```./MessagingSocketServer logLevel=DEBUG clientListenerPort=8080 eventListenerPort=3333``` <br />
Any configurations not set by commandline will default to configurations set in `conf.json`. <br />

Please use environment variables to set configurations for `./followermaz.sh` program.

## Logging
This implementation includes a custom logger. Options for logging level can be set in `conf.json`.<br />
Options include "All", "Debug", "Info", "Warn", and "Error". <br />
The default configuration is set to "INFO" but can be set to "Debug" for more in-depth look at the program.
