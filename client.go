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
	workerCount = flag.Int("workers", 100, "Amount of worker clients")
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

func startClient(wsUrl, name string, sync chan int) {
	ws, _, err := wsDialer.Dial(wsUrl, nil)
	if err != nil {
		log.Println("ERROR: ", err)
	}

	// close and signal about exiting
	defer func() {
		ws.Close()
		sync <- 0
	}()

	message := []byte("Hello!")

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
	for n := 1; n < workerCount; n++ {
		go startClient(wsUrl, fmt.Sprintf("Client[%d]", n+1), sync)
	}
	// wait for all workers
	for n := 1; n < workerCount; n++ {
		<-sync
	}
}

func main() {
	flag.Parse()
	wsUrl := url.URL{Scheme: "ws", Host: *address, Path: wsRoot}
	startSimulation(wsUrl.String(), *workerCount)
}
