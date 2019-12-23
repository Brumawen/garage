package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// Room holds the information about the room being monitored
type Room struct {
	Door1Name       string    `json:"door1name"`       // Name of garage door 1
	Door1Closed     bool      `json:"door1closed"`     // Whether door 1 is closed
	Door2Name       string    `json:"door2name"`       // Name of garage door 2
	Door2Closed     bool      `json:"door2closed"`     // Whether door 2 is closed
	Temperature     float64   `json:"temp"`            // Room temperature
	LastRead        time.Time `json:"lastread"`        // Time the values were last read
	Door1StatusTime time.Time `json:"door1statustime"` // Time that Door1 status was set
	Door2StatusTime time.Time `json:"door2statustime"` // Time that Door2 status was set
}

// WriteTo serializes the entity and writes it to the http response
func (r *Room) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}
