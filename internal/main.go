package main

import (
	"fmt"

	"github.com/nilshoeller/real-time-chat-application/internal/server"
)

func main() {
	fmt.Println("Starting...")
	
	myServer := server.NewServer()
	server.Run(myServer)
}