package main

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"
)

func runServer(t *testing.T) net.Conn {
	// start the TCP Server
	go main()

	// wait for the TCP Server to start
	time.Sleep(time.Duration(1000) * time.Millisecond)

	// connect to the TCP Server
	conn, err := net.Dial("tcp", ":9090")
	if err != nil {
		t.Fatalf("couldn't connect to the server: %v", err)
	}

	return conn
}

func sendData(t *testing.T, conn net.Conn, data string) string {
	if _, err := conn.Write([]byte(data)); err != nil {
		t.Fatalf("couldn't send request: %v", err)
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadBytes(byte('\n'))

	if err != nil {
		t.Fatalf("couldn't read server response: %v", err)
	}

	return string(response)
}

func TestWrite(t *testing.T) {
	conn := runServer(t)
	defer conn.Close()

	response := sendData(t, conn, "w000Hello\n")

	if response != "v000Hello\n" {
		t.Fatalf("unexpected server response: %s", response)
	}
}

func TestRead(t *testing.T) {
	conn := runServer(t)
	defer conn.Close()

	// First we store a value
	sendData(t, conn, "w000HelloAgain\n")

	// Then, we read that value
	response := sendData(t, conn, "r000\n")

	if response != "v000HelloAgain\n" {
		t.Fatalf("unexpected server response: %s", response)
	}
}

func TestUnsupportedCommand(t *testing.T) {
	conn := runServer(t)
	defer conn.Close()

	response := sendData(t, conn, "a000Hello\n")
	if !strings.HasPrefix(response, "e") {
		t.Fatalf("unexpected server response: %s", response)
	}
}

func TestMessageTooLong(t *testing.T) {
	conn := runServer(t)
	defer conn.Close()

	response := sendData(t, conn, "w000HelloWithAReallyLongMessageThatShouldBeWrong\n")

	if !strings.HasPrefix(response, "e") {
		t.Fatalf("unexpected server response: %s", response)
	}
}

func TestMessageTooShort(t *testing.T) {
	conn := runServer(t)
	defer conn.Close()

	response := sendData(t, conn, "r0\n")
	if !strings.HasPrefix(response, "e") {
		t.Fatalf("unexpected server response: %s", response)
	}
}

func TestMessageNotTerminated(t *testing.T) {
	conn := runServer(t)
	defer conn.Close()

	response := sendData(t, conn, "r0")
	if !strings.HasPrefix(response, "e") {
		t.Fatalf("unexpected server response: %s", response)
	}
}

func TestReadInAnotherSlot(t *testing.T) {
	conn := runServer(t)
	defer conn.Close()

	response := sendData(t, conn, "r123")
	if !strings.HasPrefix(response, "e") {
		t.Fatalf("unexpected server response: %s", response)
	}
}
