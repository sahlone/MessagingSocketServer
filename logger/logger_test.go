package logger_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sahilahmadlone/MessagingSocketServer/logger"
)

func TestServer_LoggerLevelSuccess(t *testing.T) {
	err := logger.SetLevel("all")
	if err != nil {
		t.Error(err)
	}
	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	logger.Debug("message")
	w.Close()
	output, _ := ioutil.ReadAll(r)
	os.Stdout = stdOut
	if !strings.Contains(string(output), "DEBUG") {
		t.Error("Logger level failed to set ", string(output))
	}

}
func TestServer_LoggerLevelFail(t *testing.T) {
	err := logger.SetLevel("Error")
	if err != nil {
		t.Error(err)
	}
	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	logger.Debug("Shouldn't See This")
	w.Close()
	output, _ := ioutil.ReadAll(r)
	os.Stdout = stdOut
	if strings.Contains(string(output), "DEBUG") {
		t.Error("Logger level failed to set ", string(output))
	}

}
func TestServer_LoggerLevelInvalid(t *testing.T) {
	err := logger.SetLevel("eva")
	if err == nil {
		t.Error("Incorrect Logging Level Set, should have errored")
	}

	if !strings.Contains(err.Error(), "INVALID") {
		t.Error("Should have given Invalid Log Level Error")
	}

}

func TestServer_LoggerTimeFormat(t *testing.T) {
	logger.SetTimeFormat("2006-01-02")
	logger.SetLevel("INFO")
	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	logger.Info("Testing Date Format")
	w.Close()
	output, _ := ioutil.ReadAll(r)
	os.Stdout = stdOut
	if !strings.Contains(string(output), time.Now().Format("2006-01-02")) {
		t.Error("Logger Time Format not successful ", string(output))
	}

}
