package main

import (
	"fmt"
	"time"

	telegram "github.com/brumawen/telegram/src"
)

// NotifyService handles the notifications of a door left open
type NotifyService struct {
	WasDoor1Open bool    // Indicates that Door 1 was notified as open
	WasDoor2Open bool    // Indeicates that Door 2 was notified as open
	Srv          *Server // Server
}

// Run is called from the scheduler (ClockWerk). This function will if a door has been open
// for longer than the maximum amount of time and signal an alarm
func (n *NotifyService) Run() {
	config := n.Srv.Config
	if !config.EnableDoorAlarm {
		return
	}
	if config.DoorAlarmPeriod <= 0 {
		config.DoorAlarmPeriod = 5
	}

	n.logDebug("Checking for open doors")

	room := n.Srv.Room

	// Check how long Door1 has been open
	if n.Srv.Config.EnableDoor1 {
		if room.Door1Closed {
			n.logDebug("Door1 is closed")
			if n.WasDoor1Open {
				if err := n.sendMessage(fmt.Sprintf("%s's door is now closed.", room.Door1Name)); err != nil {
					n.logError("Error notifying that door 1 is now closed.", err.Error())
				}
				n.WasDoor1Open = false
			}
		} else {
			dur := time.Since(room.Door1StatusTime)
			n.logDebug("Door 1 open for ", int(dur.Minutes()))
			if dur.Minutes() >= float64(config.DoorAlarmPeriod) {
				if err := n.sendMessage(fmt.Sprintf("%s's door has been open for %d minutes.", room.Door1Name, int(dur.Minutes()))); err != nil {
					n.logError("Error notifying that door 1 is open.", err.Error())
				}
				n.WasDoor1Open = true
			}
		}
	} else {
		n.logDebug("Door1 is disabled")
	}

	// Check how long Door2 has been open
	if n.Srv.Config.EnableDoor2 {
		if room.Door2Closed {
			n.logDebug("Door2 is closed")
			if n.WasDoor2Open {
				if err := n.sendMessage(fmt.Sprintf("%s's door is now closed.", room.Door2Name)); err != nil {
					n.logError("Error notifying that door 2 is now closed.", err.Error())
				}
				n.WasDoor2Open = false
			}
		} else {
			dur := time.Since(room.Door2StatusTime)
			n.logDebug("Door 2 open for ", int(dur.Minutes()))
			if dur.Minutes() >= float64(config.DoorAlarmPeriod) {
				if err := n.sendMessage(fmt.Sprintf("%s's door has been open for %d minutes.", room.Door2Name, int(dur.Minutes()))); err != nil {
					n.logError("Error notifying that door 2 is open.", err.Error())
				}
				n.WasDoor2Open = true
			}
		}
	} else {
		n.logDebug("Door2 is disabled")
	}
}

// sendMessage sends the specified message to telegram
func (n *NotifyService) sendMessage(m string) error {
	c := telegram.Client{
		Logger:         logger,
		VerboseLogging: n.Srv.VerboseLogging,
	}

	return c.SendMessage(m)
}

// logDebug logs a debug message to the logger
func (n *NotifyService) logDebug(v ...interface{}) {
	if n.Srv.VerboseLogging {
		a := fmt.Sprint(v...)
		logger.Info("Notify: [Dbg] ", a)
	}
}

// logError logs an error message to the logger
func (n *NotifyService) logError(v ...interface{}) {
	a := fmt.Sprint(v...)
	logger.Error("Notify: [Err] ", a)
}
