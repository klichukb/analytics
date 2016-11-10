package main

type Event struct {
	EventType string
	TS        int
	Params    map[string]interface{}
}

type Analytics int
