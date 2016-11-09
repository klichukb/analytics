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
	log.Println("serving connection")
	// Upgrader also checks this while attempting to upgrade, but in order
	// to be independent from it's implementation details, we check explicitly.
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	// Use third param for custom headers: Set-Cookie/Set-Websocket-Protocol
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	handleConnection(ws)
}

func main() {
	flag.Parse()

	http.HandleFunc(WS_ROOT, serveWebsocket)

	err := http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
