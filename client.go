package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	// "math/rand"
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

func startClient(wsUrl, name string) {
	ws, resp, err := wsDialer.Dial(wsUrl, nil)
	if err != nil {
		log.Println("ERROR: ", err)
	}
	defer ws.Close()
	log.Println(resp)

	message := []byte("Hello!")

	for {
		if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("ERROR: ", err)
			// in case this error means that servers goes down, - quit the loop.
			if websocket.IsUnexpectedCloseError(err, fatalCodes...) {
				break
			}
		}

		time.Sleep(5 * time.Second)
		// time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func startSimulation(wsUrl string, workerCount int) {
	for n := 1; n < workerCount-1; n++ {
		go startClient(wsUrl, fmt.Sprintf("Client[%d]", n+1))
	}
	startClient(wsUrl, fmt.Sprintf("Client[%d]", workerCount))
}

func main() {
	flag.Parse()
	wsUrl := url.URL{Scheme: "ws", Host: *address, Path: wsRoot}
	startSimulation(wsUrl.String(), *workerCount)
}
