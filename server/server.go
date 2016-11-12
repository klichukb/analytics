// Package provides functonality for websocket+JSON-RPC server for aggregating
// event data.
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

// Connection object, should be initialized manually.
// For convenience - use `InitDatabase` function.
var DB *sql.DB

const (
	// Root URL that websocket is accessed via.
	WsRoot = "/ws"
	// Maximum size of the message that can be read.
	readLimit = 4096
	// Time to wait for write to complete before it's considered timed out.
	writeWait = 15 * time.Second
	// Maximum time awaited for a pong back.
	pongWait = 120 * time.Second
	// Frequency of pings
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
			// Tick: time to send a ping
			// Limit the following write with a deadline.
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("PING ERROR:", err)
				return
			}
			log.Println("[Ping]")
		case <-closing:
			// If this case worked, either data has been written to the channel,
			// or it has been closed - nothing to serve, stop sending pings.
			return
		}
	}
}

// Start infinite listen loop to a websocket connection.
// Reads incoming messages, does not respond in order to spare traffic.
func HandleConnection(ws *websocket.Conn, analytics *Analytics) {
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
	// pinger will receive this signal in order to exit it's gouroutine
	defer close(closing)

	// start pinging client periodically
	go StartPinging(ws, closing)

	wrapper := &shared.WebSocketWrapper{WS: ws}
	server := rpc.NewServer()

	// register API
	server.Register(analytics)

	// start serving JSON-RPC
	codec := jsonrpc.NewServerCodec(wrapper)
	server.ServeCodec(codec)
}

// Handle HTTP request: upgrade it to websocket by replying with two headers
// Upgrade: WebSocket
// Connection: Upgrade
// Listens infinitely for new messages.
// Allowed methods: GET.
func HandleRequest(analytics *Analytics, w http.ResponseWriter, r *http.Request) {
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
	// begin the serving fun.
	HandleConnection(ws, analytics)
}

// A convenience function to initialize database connection
func InitDatabase() {
	if DB != nil {
		return
	}
	DB = GetDatabase()
}
