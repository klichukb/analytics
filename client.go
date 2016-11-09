package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/url"
	"time"
)

// Flags
var (
	address     = flag.String("address", ":8000", "Websocket server address")
	workerCount = flag.Int("workers", 1, "Amount of worker clients")
	fatalCodes  = []int{
		websocket.CloseGoingAway,
		websocket.CloseMandatoryExtension,
		websocket.CloseAbnormalClosure,
	}
)

const (
	wsRoot = "/ws"
)

var wsDialer = websocket.Dialer{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

// This function runs in seprate goroutine to ensure we ever receive pings.
// Without a reading goroutine, pings do not arrive.
func readData(ws *websocket.Conn, closing chan struct{}) {
	for {
		// a watchman channel signaling that current client is not served
		// anymore, which means that we should quit the loop.
		select {
		case <-closing:
			return
		default:
		}
		ws.ReadMessage()
	}
}

func startClient(wsUrl, name string, sync chan int) {
	ws, _, err := wsDialer.Dial(wsUrl, nil)
	if err != nil {
		log.Println("ERROR: ", err)
	}

	closing := make(chan struct{})
	// close and signal about exiting
	defer func() {
		ws.Close()
		sync <- 0
		close(closing)
	}()

	message := []byte("Hello!")
	go readData(ws, closing)

	for {
		if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("ERROR: ", err)
			return
		}
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

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
