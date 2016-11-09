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

func startSimulation(wsUrl string, workerCount int) {
	// TODO
}

func main() {
	flag.Parse()
	wsUrl := url.URL{Scheme: "ws", Host: *address, Path: wsRoot}
	startSimulation(wsUrl.String(), *workerCount)
}
