package main

import (
	"encoding/json"
	"net/http"
)

// Room holds the information about the room being monitored
type Room struct {
	Door1Name   string  `json:"door1name"` // Name of garage door 1
	Door1Open   bool    `json:"door1open"` // Whether door 1 is open
	Door2Name   string  `json:"door2name"` // Name of garage door 2
	Door2Open   bool    `json:"door2open"` // Whether door 2 is open
	Temperature float64 `json:"temp"`      // Room temperature
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
