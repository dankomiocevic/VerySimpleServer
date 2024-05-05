package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

func main() {
	l, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		return
	}

	defer l.Close()
	slot := memory_slot{value: ""}
	var mu sync.Mutex

	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}

		go handleUserConnection(conn, &slot, &mu)
	}
}

func handleUserConnection(c net.Conn, slot *memory_slot, mu *sync.Mutex) {
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
			slot.write(msg.value)
			mu.Unlock()
		}

		var sb strings.Builder
		sb.WriteString("v")
		sb.WriteString(fmt.Sprintf("%03d", msg.slot))
		sb.WriteString(slot.read())
		sb.WriteString("\n")
		c.Write([]byte(sb.String()))
	}
}
