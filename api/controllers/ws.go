package controllers

import (
	"github.com/gorilla/websocket"
	"time"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 10 * time.Second,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
}
