package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// Flags
var (
	address    = flag.String("address", ":8000", "Websocket server address")
	fatalCodes = []int{
		websocket.CloseGoingAway,
		websocket.CloseMandatoryExtension,
		websocket.CloseAbnormalClosure,
	}
)

const (
	wsRoot    = "/ws"
	readLimit = 4096
	pongWait  = 120 * time.Second
	// twice as small as time to wait for a pong back
	pingPeriod = pongWait / 2
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func updateReadDeadline(ws *websocket.Conn) {
	ws.SetReadDeadline(time.Now().Add(pongWait))
}

func handleConnection(ws *websocket.Conn) {
	defer ws.Close()

	ws.SetReadLimit(readLimit)
	updateReadDeadline(ws)

	ws.SetPongHandler(func(appData string) error {
		// update read deadline after pong
		updateReadDeadline(ws)
		return nil
	})

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("ERROR: ", err)
			// in case this error means that client goes down or leaves, we stop
			// serving it, otherwise just continue it never happened.
			if websocket.IsUnexpectedCloseError(err, fatalCodes...) {
				break
			}
		}
		log.Printf("MSG: [%v], type = %v\n", len(msg), msgType)
	}
}

func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")
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

	http.HandleFunc(wsRoot, serveWebsocket)

	log.Println("Serving...")
	err := http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
