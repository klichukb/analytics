// Package holds logic of test clients for websocket+JSON-RPC communication.
package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/klichukb/analytics/shared"
	"log"
	"math/rand"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

const (
	// Root URL that websocket is served on.
	WsRoot = "/ws"
	// Name of procedure to execute via RPC.
	proName = "Analytics.TrackEvent"
)

var wsDialer = websocket.Dialer{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

// Infinitely spamming server.
var Iterations int = -1

// Creates a sample event with random event type and mock params data.
func GenerateEvent() *shared.Event {
	eventType := shared.EventTypes[rand.Intn(len(shared.EventTypes))]
	params := map[string]interface{}{
		"var_a": 123,
		"var_b": "Foobar",
		"var_c": []int{42, 42, 84, 1, 0, 1},
	}
	return &shared.Event{eventType, int(time.Now().Unix()), params}
}

// Creates a websocket on `wsUrl` URL.
// Start single client message loop.
func StartClient(wsUrl, name string, sync chan int) {
	ws, _, err := wsDialer.Dial(wsUrl, nil)
	if err != nil {
		log.Println("ERROR: ", err)
	}

	defer func() {
		ws.Close()
		sync <- 0
	}()

	wrapper := &shared.WebSocketWrapper{WS: ws}
	codec := jsonrpc.NewClientCodec(wrapper)
	conn := rpc.NewClientWithCodec(codec)

	var reply int
	var event *shared.Event
	// We count iterations down if they are > 0.
	// Zero iterations will exit function, right away,
	// -1 iterations will spin forever, which is default value.
	iters := Iterations

	for {
		// reduce is > 0, -1 keeps spinning.
		if iters == 0 {
			break
		} else if iters > 0 {
			iters--
		}

		event = GenerateEvent()
		// sends event to server
		err = conn.Call(proName, event, &reply)
		if err != nil {
			log.Println("RPC Error: ", err)
		}
		log.Printf("Sent %v\n", event.EventType)
		time.Sleep(250 * time.Millisecond)
	}
}

// Connects to websocket on `wsUrl` URL.
// Launches `workerCount` clients that spam server with messages.
func StartSimulation(wsUrl string, workerCount int) {
	sync := make(chan int)
	for n := 0; n < workerCount; n++ {
		go StartClient(wsUrl, fmt.Sprintf("Client[%d]", n+1), sync)
	}
	// wait for all workers
	for n := 0; n < workerCount; n++ {
		<-sync
	}
}
