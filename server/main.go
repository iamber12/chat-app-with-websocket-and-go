package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	From    string          `json: "from"`
	Message string          `json: "message"`
	Sender  *websocket.Conn `json:"-"`
}

var (
	// HTTP connection upgrader
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allows connections from any Origin
		},
	}

	clients = make(map[*websocket.Conn]bool)

	msgChannel = make(chan Message)
)

// handleWsConnection upgrades HTTP connections to WebSocket and handles incoming messages.
func handlWsConnection(w http.ResponseWriter, req *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// Register the new connection
	clients[conn] = true

	for {
		// Read messages from the connection
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			conn.Close()
			delete(clients, conn)
			break
		}

		// Set the sender and send the message to the message channel
		msg.Sender = conn

		msgChannel <- msg
	}
}

// handleMessage listens on the msgChannel and broadcasts messages to all connected clients.
func handleMessage() {
	// go routine to handle incoming message
	for {
		msg := <-msgChannel

		// Broadcast the message to all clients
		for client := range clients {
			// Handle errors that occur while writing messages
			if err := client.WriteJSON(msg); err != nil {
				fmt.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// Serve client
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Configure wesocket connection
	http.HandleFunc("/ws", handlWsConnection)

	// Handle messages in websocket
	go handleMessage()

	// Setup tcp server
	http.ListenAndServe(":8000", nil)
}
