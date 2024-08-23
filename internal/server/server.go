package server

import (
	"fmt"

	"golang.org/x/net/websocket"
)

type Server struct {
	// Registered clients.
	connections map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		connections:    make(map[*websocket.Conn]bool),
	}
}

func Run(s *Server) {
	fmt.Println("Running...")
	
}