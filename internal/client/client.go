package client

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

type Client struct {
	client_name string

	server_url string
}

func NewClient(client_name string, server_url string) *Client {
	return &Client{
		client_name: client_name,
		server_url: server_url,
	}
}


func (c *Client) Run(){

	ws, err := websocket.Dial(c.server_url, "", c.client_name)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer ws.Close()

	fmt.Println("Client connected to the server.")

	// Send a message to the server
	message := "Hello Server!"
	_, err = ws.Write([]byte(message))
	if err != nil {
		log.Fatal("Write error:", err)
	}
	fmt.Println("Message sent to the server:", message)

	// Buffer for incoming messages
	buff := make([]byte, 1024)

	// Read the server's response
	n, err := ws.Read(buff)
	if err != nil {
		log.Fatal("Read error:", err)
	}
	fmt.Println("Received from server:", string(buff[:n]))

}