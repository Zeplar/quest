package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	id       int
	hub      *Hub
	socket   *websocket.Conn
	outbound chan []byte
}

func newClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		id:       len(hub.clients),
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
}

func (client *Client) read() {
	defer func() {
		client.hub.unregister <- client
	}()
	for {
		_, data, err := client.socket.ReadMessage()
		if err != nil {
			break
		}
		client.hub.onMessage(data, client)
	}
}

func (client *Client) write() {
	for data := range client.outbound {
		client.socket.WriteMessage(websocket.TextMessage, data)
	}
	client.socket.WriteMessage(websocket.CloseMessage, []byte{})
}

func (client Client) run() {
	go client.read()
	go client.write()
}

func (client Client) close() {
	client.socket.Close()
	close(client.outbound)
}
