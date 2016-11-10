package main

import (
	"flag"
	"github.com/klichukb/analytics/client"
	"github.com/klichukb/analytics/server"
	"log"
	"net/http"
	"net/url"
)

var (
	mode        = flag.String("mode", "server", "Use `server` or `client`")
	address     = flag.String("address", ":8000", "Websocket address")
	workerCount = flag.Int("workers", 100, "Amount of worker clients")
)

func runServer() {
	server.InitDatabase()
	defer server.DB.Close()

	http.HandleFunc(server.WsRoot, server.HandleRequest)
	log.Println("Serving...")

	err := http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}

func runClient() {
	wsUrl := url.URL{Scheme: "ws", Host: *address, Path: client.WsRoot}
	client.StartSimulation(wsUrl.String(), *workerCount)
}

var modes = map[string]func(){
	"server": runServer,
	"client": runClient,
}

func main() {
	flag.Parse()

	handler := modes[*mode]
	handler()
}
