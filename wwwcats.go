package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/eyedeekay/sam3/helper"
	"github.com/eyedeekay/sam3/i2pkeys"
)

var addr = flag.String("l", ":8080", "http service address")

var REVISION = 9

func main() {
	flag.Parse()

	// Create a global list of lobbies
	lobbies := make(map[string]*Lobby)

	// Serve the client-side software
	fs := http.FileServer(http.Dir("public_html"))
	http.Handle("/", fs)

	// Handle incoming websocket connections
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(w, r, lobbies)
	})

	// Start the server

	i2plistener, err := sam.I2PListener("wwwcats","127.0.0.1:7656","wwwcats")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Now listening on", "http://"+i2plistener.Addr().(i2pkeys.I2PAddr).Base32())
	log.Fatal(http.Serve(i2plistener, nil))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Upgrade incoming connections to websockets
func handleConnections(w http.ResponseWriter, r *http.Request, lobbies map[string]*Lobby) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Instantiate the new client object
	client := &Client{conn: conn, send: make(chan []byte, 256)}

	// Hand the client off to these goroutines which will handle all i/o
	go client.readPump(lobbies)
	go client.writePump()
}
