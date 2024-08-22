package chat

import "fmt"

type Server struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
	}
}

func Run(s *Server) {
	fmt.Println("Running...")
	for {
		select {
		case message := <-s.broadcast:
			for client := range s.clients {
				fmt.Println("client: ", client)
				fmt.Println("message: ", message)
		}
	}
}
}