package main

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func TestEcho(t *testing.T) {
	// start the TCP Server
	go main()

	// wait for the TCP Server to start
	time.Sleep(1 * time.Second)

	// connect to the TCP Server
	conn, err := net.Dial("tcp", ":9090")
	if err != nil {
		t.Fatalf("couldn't connect to the server: %v", err)
	}
	defer conn.Close()

	// test the TCP Server output
	if _, err := conn.Write([]byte("END\n")); err != nil {
		t.Fatalf("couldn't send request: %v", err)
	} else {
		reader := bufio.NewReader(conn)
		response, err := reader.ReadBytes(byte('\n'))

		if err != nil {
			t.Fatalf("couldn't read server response: %v", err)
		}

		if string(response) != "END\n" {
			t.Fatalf("unexpected server response: %s", response)
		}
	}
}
