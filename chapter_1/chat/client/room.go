package client

import (
	"log"
	"net/http"

	"github.com/oz117/go_blueprints/chapter_1/chat/trace"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: messageBufferSize}

// Room represents a chat room
type Room struct {
	// Channel that contains the data to be forwarded to all the clients
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	// To help debug and understand what's going on inside the room
	Tracer trace.Tracer
}

// NewRoom creates a room
func NewRoom() *Room {
	return &Room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		Tracer:  trace.Off(),
	}
}

// Run watches the three room channels to see
// if a client joins the room
// if a client leaves the room
// if there is any messege that can be sent to the clients of the room
func (r *Room) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.Tracer.Trace("A new client has joined the room")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.Tracer.Trace("A client has left the room")
		case msg := <-r.forward:
			r.Tracer.Trace("Message received: [", string(msg), "]")
			// Forward a message to all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.Tracer.Trace(" -- sent to client")
				default:
					delete(r.clients, client)
					close(client.send)
					r.Tracer.Trace(" -- failed to send, cleaned up client")
				}
			}
		}
	}
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
