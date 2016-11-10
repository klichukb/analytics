package server

import (
	"database/sql"
	"github.com/gorilla/websocket"
	"github.com/klichukb/analytics/shared"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

var DB *sql.DB

const (
	WsRoot = "/ws"

	readLimit = 4096
	// time to wait for write to complete
	writeWait = 15 * time.Second
	pongWait  = 120 * time.Second
	// twice as small as time to wait for a pong back
	pingPeriod = pongWait / 2
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

// Sets read deadline to `now` + `pongWait`.
func UpdateReadDeadline(ws *websocket.Conn) {
	ws.SetReadDeadline(time.Now().Add(pongWait))
}

// Launch a loop of pings, based on timer.
// Will however obey to `closing` channel and stop the loop
// whenver channel gets closed from outside.
func StartPinging(ws *websocket.Conn, closing chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("PING ERROR:", err)
				return
			}
			log.Println("[Ping]")
		case <-closing:
			return
		}
	}
}

// Start infinite listen loop to a websocket connection.
// Reads incoming messages, does not respond in order to spare traffic.
func HandleConnection(ws *websocket.Conn) {
	ws.SetReadLimit(readLimit)
	UpdateReadDeadline(ws)

	// make sure we know how to handle response.
	ws.SetPongHandler(func(appData string) error {
		log.Println("[Pong]")
		// update read deadline after pong
		UpdateReadDeadline(ws)
		return nil
	})

	closing := make(chan struct{})
	defer close(closing)
	go StartPinging(ws, closing)

	wrapper := &shared.WebSocketWrapper{WS: ws}
	server := rpc.NewServer()

	// register API
	server.Register(new(Analytics))

	codec := jsonrpc.NewServerCodec(wrapper)
	server.ServeCodec(codec)
}

// Handle HTTP request: upgrade it to websocket by replying with two headers
// Upgrade: WebSocket
// Connection: Upgrade
// Listens infinitely for new messages.
// Allowed methods: GET.
func HandleRequest(w http.ResponseWriter, r *http.Request) {
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

	defer ws.Close()
	HandleConnection(ws)
}

func InitDatabase() {
	DB = GetDatabase()
}
