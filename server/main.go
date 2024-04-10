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

func handlWsConnection(w http.ResponseWriter, req *http.Request) {
	// Update http connection to websocket
	conn, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	clients[conn] = true

	for {
		// handle message
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			break
		}

		msg.Sender = conn

		msgChannel <- msg
	}
}

func handleMessage() {
	// go routine to handle incoming message
	for {
		msg := <-msgChannel

		// broadcast message to all users
		for client := range clients {
			if err := client.WriteJSON(msg); err != nil {
				fmt.Printf("error: %v", err)
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
