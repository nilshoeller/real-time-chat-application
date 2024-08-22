package main

import (
	"fmt"

	"github.com/nilshoeller/real-time-chat-application/internal/chat"
)

func main() {
	fmt.Println("Starting...")
	
	server := chat.NewServer()
	chat.Run(server)
}