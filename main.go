package main

import (
	"log"

	"github.com/arsidada/gas-price-bot/server"
)

func main() {
	err := server.StartServer()
	if err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}
