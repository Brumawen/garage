package main

import (
	"fmt"
	"time"
)

// Thingspeak uploads the room telemetry to the Thingspeak server in the cloud
type Thingspeak struct {
	Srv        *Server   // Current Server
	LastUpdate time.Time // Last time the update was run
	IsRunning  bool      // Is this currently running?
}

// Run is called from the scheduler (ClockWerk). This function will get the latest measurements
// and send the measurements to Thingspeak
func (t *Thingspeak) Run() {
	t.logInfo("Uploading telemetry")
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
