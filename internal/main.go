package main

import (
	"fmt"

	"github.com/nilshoeller/real-time-chat-application/internal/server"
)

func main() {
	fmt.Println("Starting...")
	
	server.Run()
}