package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/zeplar/quest/message"
)

type Hub struct {
	clients              []*Client
	register, unregister chan *Client
	strokes              []message.Stroke
}

func newHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

var upgrader = websocket.Upgrader{
	// Allow all origins
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (hub *Hub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not upgrade", http.StatusInternalServerError)
		return
	}
	client := newClient(hub, socket)
	hub.clients = append(hub.clients, client)
	hub.register <- client
	client.run()
}

func (hub *Hub) send(message interface{}, client *Client) {
	data, _ := json.Marshal(message)
	client.outbound <- data
}

func (hub *Hub) broadcast(message interface{}, ignore *Client) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	for _, c := range hub.clients {
		if c != ignore {
			c.outbound <- data
		}
	}
}

func (hub *Hub) onConnect(client *Client) {
	log.Println("client connected: ", client.socket.RemoteAddr())
	hub.initialize(client)
}

func (hub *Hub) initialize(client *Client) {
	msg := message.Message{Kind: message.KindStroke}
	for i := range hub.strokes {
		msg.Stroke = hub.strokes[i]
		hub.send(msg, client)
	}
	msg = message.Message{Kind: message.KindConnected}
	msg.Stroke.OwnerID = len(hub.clients)
	hub.send(msg, client)
}

func (hub *Hub) onDisconnect(client *Client) {
	log.Println("client disconnected: ", client.socket.RemoteAddr())
	client.close()
	// Find index of client
	i := -1
	for j, c := range hub.clients {
		if c.id == client.id {
			i = j
			break
		}
	}
	// Delete client from list
	hub.clients = append(hub.clients[:i], hub.clients[i+1:]...)
}

func (hub *Hub) onMessage(data []byte, client *Client) {
	var msg message.Message
	log.Println(string(data))
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Println(err)
		return
	}
	switch msg.Kind {
	case message.KindStroke:
		hub.strokes = append(hub.strokes, msg.Stroke)
	case message.KindUndo:
	case message.KindClear:
	}
	hub.broadcast(msg, client)
}
