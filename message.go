package main

import (
	"errors"
	"strconv"
	"strings"
)

type message struct {
	command byte
	slot    int
	value   string
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
