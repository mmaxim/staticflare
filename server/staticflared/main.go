package main

import (
	"log"

	"mmaxim.org/staticflare/server"
)

func main() {
	if err := server.NewServer("localhost:8080").Run(); err != nil {
		log.Printf("error running server: %s\n", err)
	}
}
