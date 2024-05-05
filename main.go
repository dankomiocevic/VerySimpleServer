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
	slots := configureSlots()

	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}

		go handleUserConnection(conn, slots)
	}
}

func configureSlots() [1000]Slot {
	slot_0 := memory_slot{value: ""}
	slot_1 := memory_slot{value: ""}
	slots := [1000]Slot{}
	slots[0] = &slot_0
	slots[1] = &slot_1

	return slots
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

		if msg.command == 'w' {
			current_slot.write(msg.value)
		}

		var sb strings.Builder
		sb.WriteString("v")
		sb.WriteString(fmt.Sprintf("%03d", msg.slot))
		sb.WriteString(current_slot.read())
		sb.WriteString("\n")
		c.Write([]byte(sb.String()))
	}
}
