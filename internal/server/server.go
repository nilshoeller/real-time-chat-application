package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	// Registered clients.
	connections map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		connections: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())

	// not concurrent save -> use mutex
	s.connections[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buff := make([]byte, 1024)
	for {
		n, err := ws.Read(buff)
		if err != nil {
			if err == io.EOF {
				// delete(s.connections, ws) // needed ???
				break
			}
			fmt.Println("read error:", err)
			continue
		}
		msg := buff[:n]
		fmt.Println(string(msg))
		ws.Write([]byte(ws.RemoteAddr().String() + ": " + string(msg)))
		// ws.Write([]byte("thank you for the message!"))
	}
}

func (s *Server) Run() {
	fmt.Println("Running on :3000")

	http.Handle("/ws", websocket.Handler(s.handleWS))
	http.ListenAndServe(":3000", nil)

}

func main() {
	myServer := NewServer()
	myServer.Run()
}
