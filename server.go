package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

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
	http.HandleFunc(WS_ROOT, serveWebsocket)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
