package server_test

import (
	"testing"
	"os"
	"os/exec"
	"net"
	"io"
	"strings"

	"github.com/sahilahmadlone/MessagingSocketServer/logger"
	"github.com/sahilahmadlone/MessagingSocketServer/config"
	"github.com/sahilahmadlone/MessagingSocketServer/server"
)

func TestServer_StartAndStop(t *testing.T) {
	conf := config.ServerDefaultConfig("../config/")
	logger.SetLevel("ERROR")
	s, err := server.Run(*conf)
	if err != nil {
		recover()
		t.Error("Server didn't start properly ", err)
	}
	err = s.ShutDown()
	if err != nil {
		t.Error("Server didn't shut down properly", err)
	}
	if s.IsRunning {
		t.Error("Server still running! ", s.IsRunning)
	}
}
func TestServer_Run(t *testing.T) {
	conf := config.ServerDefaultConfig("../config/")
	s, err := server.Run(*conf)
	if err != nil {
		t.Error("Error starting server")
		return
	}
	os.Setenv("totalEvents", "1000")
	out, err := exec.Command("time", "java", "-server", "-Xmx1G", "-jar", "../follower-maze-2.0.jar").Output()
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(string(out), "ALL NOTIFICATIONS RECEIVED") {
		t.Error("Failed test run (against jar file)")
	}
	s.ShutDown()
}

func TestServer_RunWithBadEvent(t *testing.T) {
	logger.SetLevel("ERROR")
	conf := config.ServerDefaultConfig("../config/")
	s, err := server.Run(*conf)
	if err != nil {
		t.Error("Error starting server ", err)

	}
	badEvents := []string{"sldjfs", "3248nfk", "1|1|1|1|1|1", ""}
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		t.Error(err)

	}

	for _, e := range badEvents {
		_, err := io.WriteString(conn, e)
		if err != nil {
			t.Error(err)

		}
	}

	conn.Close()
	if err := s.ShutDown(); err != nil {
		t.Error(err)
	}
}

func TestServer_RunWithBadUsers(t *testing.T) {
	logger.SetLevel("Error")
	conf := config.ServerDefaultConfig("../config/")
	s, err := server.Run(*conf)
	if err != nil {
		t.Error("Error starting server ", err)
		s.ShutDown()

	}
	badUsers := []string{"sldjfs", "3248nfk", "1|1|1|1|1|1", ""}
	conn, err := net.Dial("tcp", "localhost:9099")
	if err != nil {
		t.Error(err)
		s.ShutDown()

	}

	for _, e := range badUsers {
		_, err := io.WriteString(conn, e)
		if err != nil {
			t.Error(err)

		}
	}
	conn.Close()
	s.ShutDown()
}

//Benchmarking tests
func benchmarkServer(b *testing.B, numEvents string) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		conf := config.ServerDefaultConfig("../config/")
		s, err := server.Run(*conf)
		if err != nil {
			b.Error("Error starting server")

		}
		os.Setenv("totalEvents", numEvents)
		out, err := exec.Command("time", "java", "-server", "-Xmx1G", "-jar", "../follower-maze-2.0.jar").Output()
		if err != nil {
			b.Error(err)
		}
		if !strings.Contains(string(out), "ALL NOTIFICATIONS RECEIVED") {
			b.Error("Failed test run (against jar file)")
		}
		s.ShutDown()
	}

}

func Benchmark1ThousandEvents(b *testing.B)         { benchmarkServer(b, "1000") }
func BenchmarkServer10ThousandEvents(b *testing.B)  { benchmarkServer(b, "10000") }
func BenchmarkServer100ThousandEvents(b *testing.B) { benchmarkServer(b, "100000") }
