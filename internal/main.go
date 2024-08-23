package main

import (
	"github.com/nilshoeller/real-time-chat-application/internal/server"
)

func main() {
	// fmt.Println("Starting...")
	
	myServer := server.NewServer()
	myServer.Run()
}