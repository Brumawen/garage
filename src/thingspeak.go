package main

import (
	"fmt"
	"net/http"
	"time"
)

// Thingspeak uploads the room telemetry to the Thingspeak server in the cloud
type Thingspeak struct {
	Srv        *Server   // Current Server
	LastUpdate time.Time // Last time the update was run
	IsRunning  bool      // Is this currently running?
	lastValues *Room     // Last values uploaded for Room
}

// Run is called from the scheduler (ClockWerk). This function will get the latest measurements
// and send the measurements to Thingspeak
func (t *Thingspeak) Run() {
	if !t.Srv.Config.EnableThingspeak {
		t.logInfo("Thingspeak has been disabled")
		return
	}
	key := t.Srv.Config.ThingspeakID
	if key == "" {
		t.logError("Thingspeak API ID has not been configured")
		return
	}

	mustUpload := false
	if t.lastValues == nil {
		t.lastValues = &Room{}
		mustUpload = true
	} else {
		// Check for changes
		if t.lastValues.Door1Closed != t.Srv.Room.Door1Closed || t.lastValues.Door2Closed != t.Srv.Room.Door2Closed || t.lastValues.Temperature != t.Srv.Room.Temperature {
			mustUpload = true
		}
	}

	if mustUpload {
		t.logInfo("Uploading telemetry")
		client := http.Client{}
		door1Closed := 0
		if t.Srv.Room.Door1Closed {
			door1Closed = 1
		}
		door2Closed := 0
		if t.Srv.Room.Door2Closed {
			door2Closed = 1
		}
		url := fmt.Sprintf("https://api.thingspeak.com/update?api_key=%s&field1=%d&field2=%d&field3=%d&field4=%f",
			key,
			door1Closed,
			door2Closed,
			1,
			t.Srv.Room.Temperature)
		if resp, err := client.Get(url); err != nil {
			t.logError("Error sending telemetry to Thingspeak.", err.Error())
		} else {
			if resp.StatusCode != 200 {
				t.logError("Error sending telemetry to Thingspeak. Status", resp.StatusCode, "returned.")
			}
		}
	}

	// Update last values
	t.lastValues.Door1Name = t.Srv.Room.Door1Name
	t.lastValues.Door1Closed = t.Srv.Room.Door1Closed
	t.lastValues.Door2Name = t.Srv.Room.Door2Name
	t.lastValues.Door2Closed = t.Srv.Room.Door2Closed
	t.lastValues.Temperature = t.Srv.Room.Temperature
	t.lastValues.LastRead = t.Srv.Room.LastRead
}

// logInfo logs an information message to the logger
func (t *Thingspeak) logInfo(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("Thingspeak: [Inf] ", a[1:len(a)-1])
}

// logError logs an error message to the logger
func (t *Thingspeak) logError(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Error("Thingspeak [Err] ", a[1:len(a)-1])
}
