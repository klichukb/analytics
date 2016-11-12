// Defines API for RPC of server.
// type Analytics struct -- a structure that servers for exposing public functionality
// to RPC API.
package server

import (
	"errors"
	"github.com/klichukb/analytics/shared"
	"log"
	"sync"
	"time"
)

// Api object that exports its methods to RPC.
type Analytics struct {
	eventBuffer []*shared.Event
	mu          sync.Mutex
	lastEvent   time.Time

	// set to true in order to kill event buffer watcher goroutine
	StopBufferWatch chan int
}

var (
	MaxBufferSize = 100
	FlushTimeout  = 500 * time.Millisecond
)

// Flush event buffer to DB in case it's not empty.
func (analytics *Analytics) flushEventBuffer() error {
	if len(analytics.eventBuffer) == 0 {
		return nil
	}
	err := SaveEvents(analytics.eventBuffer...)
	if err == nil {
		// reset buffer to empty
		analytics.eventBuffer = analytics.eventBuffer[0:0]
	}
	log.Println("Flushed events to DB")
	return err
}

// Watches for any leftover buffer to flush when there have been
// no new events coming in recently.
func (analytics *Analytics) watchEventBuffer() {
	ticker := time.NewTicker(FlushTimeout)
	for {
		select {
		case <-ticker.C:
			// wrapping into a function makes sure its releasing the mutex after it's
			// body execution
			func() {
				analytics.mu.Lock()
				defer analytics.mu.Unlock()

				// flush the event buffer if `FlushTimeout` has pased since last event coming in.
				if time.Since(analytics.lastEvent) >= FlushTimeout {
					analytics.flushEventBuffer()
				}
			}()
		case <-analytics.StopBufferWatch:
			return
		}
	}
}

// Creates a new analytics object and starts a goroutine
// that watches for any leftover buffer to flush when there have been
// no new events coming in recently.
func NewAnalytics() *Analytics {
	analytics := &Analytics{lastEvent: time.Now(), StopBufferWatch: make(chan int)}
	go analytics.watchEventBuffer()
	return analytics
}

// Save event to buffer/DB.
// Event is not guaranteed to be persisted immediatly, instead can be written
// to buffer for future bulk-wite to DB.
func (analytics *Analytics) addEvent(event *shared.Event) error {
	analytics.mu.Lock()
	defer analytics.mu.Unlock()

	if len(analytics.eventBuffer) >= MaxBufferSize {
		// flush buffer - persist to database
		return analytics.flushEventBuffer()
	} else {
		analytics.eventBuffer = append(analytics.eventBuffer, event)
		analytics.lastEvent = time.Now()
		return nil
	}
}

// Process event data.
func (analytics *Analytics) TrackEvent(event *shared.Event, reply *int) error {
	if len(event.EventType) == 0 || event.TS == 0 {
		// invalid parameters
		return errors.New("Invalid event data")
	}
	log.Println("Event received:", event, &event)

	return analytics.addEvent(event)
}
