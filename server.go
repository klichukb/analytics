package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// Flags
var address = flag.String("address", ":8000", "Websocker server address")

const WS_ROOT = "/ws"

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func handleConnection(ws *websocket.Conn) {
	// TODO
}

func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func main() {
	flag.Parse()

	http.HandleFunc(WS_ROOT, serveWebsocket)

	err := http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
