package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		return
	}

	defer l.Close()
	conf, err := loadConfig("example_config.yml")
	if err != nil {
		return
	}

	slots := configureSlots(conf)

	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}

		go handleUserConnection(conn, slots)
	}
}

func handleUserConnection(c net.Conn, slots [1000]Slot) {
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

		current_slot := slots[msg.slot]
		if current_slot == nil {
			c.Write([]byte("e\n"))
			continue
		}

		var value string
		if msg.command == 'w' {
			value, err = current_slot.write(msg.value, c)

			if err != nil {
				c.Write([]byte("e\n"))
				continue
			}
		} else {
			value = current_slot.read()
		}

		var sb strings.Builder
		sb.WriteString("v")
		sb.WriteString(fmt.Sprintf("%03d", msg.slot))
		sb.WriteString(value)
		sb.WriteString("\n")
		c.Write([]byte(sb.String()))
	}
}
