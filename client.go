package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/rpc"
	"net/rpc/jsonrpc"
	"net/url"
	"time"
)

// Flags
var (
	address     = flag.String("address", ":8000", "Websocket server address")
	workerCount = flag.Int("workers", 4, "Amount of worker clients")
	fatalCodes  = []int{
		websocket.CloseGoingAway,
		websocket.CloseMandatoryExtension,
		websocket.CloseAbnormalClosure,
	}
)

const (
	wsRoot  = "/ws"
	proName = "Analytics.TrackEvent"
)

var wsDialer = websocket.Dialer{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func generateEvent() *Event {
	params := map[string]interface{}{
		"var_a": 123,
		"var_b": "Foobar",
		"var_c": []int{42, 42, 84, 1, 0, 1},
	}
	return &Event{"PageView", int(time.Now().Unix()), params}
}

// Creates a websocket on `wsUrl` URL.
// Start single client message loop.
func startClient(wsUrl, name string, sync chan int) {
	ws, _, err := wsDialer.Dial(wsUrl, nil)
	if err != nil {
		log.Println("ERROR: ", err)
	}

	defer func() {
		ws.Close()
		sync <- 0
	}()

	wrapper := &WebSocketWrapper{ws: ws}
	codec := jsonrpc.NewClientCodec(wrapper)
	conn := rpc.NewClientWithCodec(codec)

	var reply int
	for {
		err = conn.Call(proName, generateEvent(), &reply)
		if err != nil {
			log.Println("RPC Error: ", err)
		}
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

// Connects to websocket on `wsUrl` URL.
// Launches `workerCount` clients that spam server with messages.
func startSimulation(wsUrl string, workerCount int) {
	sync := make(chan int)
	for n := 0; n < workerCount; n++ {
		go startClient(wsUrl, fmt.Sprintf("Client[%d]", n+1), sync)
	}
	// wait for all workers
	for n := 0; n < workerCount; n++ {
		<-sync
	}
}

func main() {
	flag.Parse()
	wsUrl := url.URL{Scheme: "ws", Host: *address, Path: wsRoot}
	startSimulation(wsUrl.String(), *workerCount)
}
