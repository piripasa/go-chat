package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

type Message struct {
	Email 	 string `json:"email"`
	Username string `json:"Username"`
	Message  string `json:"Message"`
}

func main()  {
	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request)  {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			// Send the newly received message to the broadcast channel
			broadcast <- msg
		}
	}
}