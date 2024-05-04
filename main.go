package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type message struct {
	command byte
	slot    int
	value   string
}

func main() {
	l, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		return
	}

	defer l.Close()
	slot := ""
	var mu sync.Mutex

	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}

		go handleUserConnection(conn, &slot, &mu)
	}
}

func handleUserConnection(c net.Conn, slot *string, mu *sync.Mutex) {
	defer func() {
		c.Close()
	}()

	buf := make([]byte, 41)

	for {
		size, err := bufio.NewReader(c).Read(buf)
		if err != nil {
			return
		}

		msg, err := ParseMessage(size, buf)
		if err != nil {
			c.Write([]byte("e\n"))
			continue
		}

		if msg.command == 'w' {
			mu.Lock()
			*slot = msg.value
			mu.Unlock()
		}

		var sb strings.Builder
		sb.WriteString("v")
		sb.WriteString(fmt.Sprintf("%03d", msg.slot))
		sb.WriteString(*slot)
		sb.WriteString("\n")
		c.Write([]byte(sb.String()))
	}
}

func ParseMessage(size int, buf []byte) (*message, error) {
	if buf[size-1] != '\n' {
		return nil, errors.New("Message is malformed")
	}
	input := strings.Split(string(buf), "\n")[0]

	if len(input) < 4 {
		return nil, errors.New("Message is too short")
	}

	if len(input) > 40 {
		return nil, errors.New("Message is too long")
	}

	command := input[:1]

	if command != "r" && command != "w" {
		return nil, errors.New("Command not supported")
	}

	slot, err := strconv.Atoi(input[1:4])
	if err != nil {
		return nil, errors.New("Malformed slot")
	}

	// Only one slot for now
	if slot > 0 {
		return nil, errors.New("Slot not supported")
	}

	var value string
	if command == "w" {
		value = input[4:]
	}

	return &message{command: []byte(command)[0], slot: slot, value: value}, nil
}
