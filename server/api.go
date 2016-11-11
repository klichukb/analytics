// Defines API for RPC of server.
// type Analytics struct -- a structure that servers for exposing public functionality
// to RPC API.
package server

import (
	"errors"
	"github.com/klichukb/analytics/shared"
	"log"
)

// Api object that exports its methods to RPC.
type Analytics struct{}

// Process/persist event data.
func (comm *Analytics) TrackEvent(event *shared.Event, reply *int) error {
	if len(event.EventType) == 0 || event.TS == 0 {
		// invalid parameters
		return errors.New("Invalid event data")
	}
	log.Println("Event received:", event, &event)

	// persist to database
	err := SaveEvent(event)
	if err != nil {
		return err
	}
	return nil
}
