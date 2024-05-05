package main

import (
	"sync"
)

type memory_slot struct {
	value string
	mu    sync.Mutex
}

type slot interface {
	read() string
	write(string) string
}

func (m *memory_slot) read() string {
	return m.value
}

func (m *memory_slot) write(data string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.value = data
}
