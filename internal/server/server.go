package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	// Registered clients.
	connections map[*websocket.Conn]bool
}

type MessageData struct {
	ClientName string `json:"clientName"`
	Message    string `json:"message"`
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
	for {
		buff := make([]byte, 1024)

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

		var msgData MessageData
		err = json.Unmarshal(msg, &msgData)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			continue
		}

		response := msgData.ClientName + ": " + msgData.Message
		fmt.Println("Received message - ", response)
		ws.Write([]byte(response))
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
