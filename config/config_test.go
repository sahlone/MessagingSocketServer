package config_test

import (
	"testing"
	"reflect"
	"github.com/sahilahmadlone/MessagingSocketServer/config"
)

func TestServerConfigShouldEqual(t *testing.T) {
	conf:=config.ServerDefaultConfig("./")
	msc := config.ServerConfig{LogLevel:"INFO", ClientListenerPort:9099, EventListenerPort:9090, SequenceNumber:1}
	if !reflect.DeepEqual(*conf, msc) {
		t.Error("Configurations are NOT equal")
	}
}
func TestServerConfigShouldNotEqual(t *testing.T) {
	conf:=config.ServerDefaultConfig("./")
	msc := config.ServerConfig{LogLevel:"INFO", ClientListenerPort:9090, EventListenerPort:9090, SequenceNumber:1}
	if reflect.DeepEqual(*conf, msc) {
		t.Error("Configurations are equal and should NOT be")
	}
}
