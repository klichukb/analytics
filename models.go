package main

type Event struct {
	EventType string                 `json:"event_type"`
	TS        int                    `json:"ts"`
	Params    map[string]interface{} `json:"params"`
}
