package client

import (
	"golang.org/x/net/websocket"

	"github.com/nilshoeller/real-time-chat-application/internal/server"
)

type Client struct {
	server *server.Server

	conn *websocket.Conn
}