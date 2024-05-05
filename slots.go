package main

type memory_slot struct {
	value string
}

type slot interface {
	read() string
	write(string) string
}

func (m *memory_slot) read() string {
	return m.value
}

func (m *memory_slot) write(data string) {
	m.value = data
}
