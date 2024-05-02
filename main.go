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

	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}

		go handleUserConnection(conn)
	}
}

func handleUserConnection(c net.Conn) {
	defer func() {
		c.Close()
	}()

	for {
		clientInput, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			return
		}
		fmt.Println(clientInput)
		c.Write([]byte(clientInput))

		if strings.TrimRight(clientInput, "\n") == "END" {
			return
		}
	}
}
