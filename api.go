package main

import "errors"
import "log"

type Analytics int

func (comm *Analytics) TrackEvent(event *Event, reply *int) error {
	if len(event.EventType) == 0 || event.TS == 0 {
		// invalid parameters
		*reply = 1
		return errors.New("Invalid event data")
	}
	log.Println("Event received:", event, &event)
	*reply = 0
	return nil
}
