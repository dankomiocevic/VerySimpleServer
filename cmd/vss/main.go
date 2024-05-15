package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/dankomiocevic/VerySimpleServer/internal/config"
	"github.com/dankomiocevic/VerySimpleServer/internal/server"
	"github.com/dankomiocevic/VerySimpleServer/internal/slots"
)

func main() {
	l, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		return
	}

	defer l.Close()
	conf, err := config.LoadConfig("../../example_config.yml")
	if err != nil {
		return
	}

	slotsArray := config.ConfigureSlots(conf)

	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}

		go handleUserConnection(conn, slotsArray)
	}
}

func handleUserConnection(c net.Conn, slotsArray [1000]slots.Slot) {
	defer func() {
		c.Close()
	}()

	buf := make([]byte, 41)

	for {
		size, err := bufio.NewReader(c).Read(buf)
		if err != nil {
			return
		}

		msg, err := server.ParseMessage(size, buf)
		if err != nil {
			c.Write([]byte("e\n"))
			continue
		}

		current_slot := slotsArray[msg.Slot]
		if current_slot == nil {
			c.Write([]byte("e\n"))
			continue
		}

		var value string
		if msg.Command == 'w' {
			value, err = current_slot.Write(msg.Value, c)

			if err != nil {
				c.Write([]byte("e\n"))
				continue
			}
		} else {
			value = current_slot.Read()
		}

		var sb strings.Builder
		sb.WriteString("v")
		sb.WriteString(fmt.Sprintf("%03d", msg.Slot))
		sb.WriteString(value)
		sb.WriteString("\n")
		c.Write([]byte(sb.String()))
	}
}
