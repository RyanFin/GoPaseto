package main

import (
	"RyanFin/GoPaseto/server"
	"log"
)

const (
	addr = "0.0.0.0:8080"
)

func main() {
	_, err := server.NewServer(addr)
	if err != nil {
		log.Fatal("Failed to start server.")
	}
}
