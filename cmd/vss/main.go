package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dankomiocevic/VerySimpleServer/internal/config"
	"github.com/dankomiocevic/VerySimpleServer/internal/server"
)

func main() {
	fmt.Println("Loading config..")
	conf, err := config.LoadConfig("../../example_config.yml")
	if err != nil {
		return
	}

	fmt.Println("Initializing slots..")
	slotsArray := config.ConfigureSlots(conf)

	fmt.Println("Starting server..")
	s := server.NewServer("localhost:9090", slotsArray)
	defer s.Stop()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	fmt.Println("Shutting down server..")
}
