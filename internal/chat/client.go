package chat

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	server *Server

	conn *websocket.Conn
}