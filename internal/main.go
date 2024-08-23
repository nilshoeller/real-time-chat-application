package main

import (
	"time"

	"github.com/nilshoeller/real-time-chat-application/internal/client"
	"github.com/nilshoeller/real-time-chat-application/internal/server"
)

func main() {
	// fmt.Println("Starting...")
	
	myServer := server.NewServer()
	go myServer.Run()

	time.Sleep(time.Second * 5)

	newClient := client.NewClient("http://this-is-a-new-client:8000/", "ws://localhost:3000/ws")
	newClient.Run()
}