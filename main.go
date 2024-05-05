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
	slot := memory_slot{value: ""}

	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}

		go handleUserConnection(conn, &slot)
	}
}

func handleUserConnection(c net.Conn, slot *memory_slot) {
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
			slot.write(msg.value)
		}

		var sb strings.Builder
		sb.WriteString("v")
		sb.WriteString(fmt.Sprintf("%03d", msg.slot))
		sb.WriteString(slot.read())
		sb.WriteString("\n")
		c.Write([]byte(sb.String()))
	}
}
