// Launches a server/client for Websockes + JSON-RPC communication.
// Works in both modes.

// Use --mode option to choose: (server|client)

// For client --workers option sets quantity of clients spinning.
// Both support --address option.

// Examples:
// Start server
//    analytics --mode server
//
// Start 10 concurrent clients.
//    analytics --mode client --workers 10

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

// Launch server on default port
func runServer() {
	server.InitDatabase()
	defer server.DB.Close()

	analytics := server.NewAnalytics()
	http.HandleFunc(server.WsRoot, func(w http.ResponseWriter, r *http.Request) {
		server.HandleRequest(analytics, w, r)
	})

	log.Println("Serving...")

	err := http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}

// Launch test clients to feed the server.
func runClient() {
	wsUrl := url.URL{Scheme: "ws", Host: *address, Path: client.WsRoot}
	// Starts `workerCount` clients.
	client.StartSimulation(wsUrl.String(), *workerCount)
}

// Maps options
var modes = map[string]func(){
	"server": runServer,
	"client": runClient,
}

func main() {
	flag.Parse()

	handler := modes[*mode]
	// Only run supported modes, otherwise fail & bye.
	if handler == nil {
		log.Fatal("Unsupported mode")
	}
	// Runs either client or server
	handler()
}
