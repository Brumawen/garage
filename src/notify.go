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

	room := n.Srv.Room

	// Check how long Door1 has been open
	if room.Door1Closed {
		if n.WasDoor1Open {
			if err := n.sendMessage(fmt.Sprintf("%s's door is now closed.", room.Door1Name)); err != nil {
				n.logError("Error notifying that door 1 is now closed.", err.Error())
			}
			n.WasDoor1Open = false
		}
	} else {
		dur := time.Since(room.Door1StatusTime)
		if dur.Minutes() >= float64(config.DoorAlarmPeriod) {
			if err := n.sendMessage(fmt.Sprintf("%s's door has been open for %d minutes.", room.Door1Name, int(dur.Minutes()))); err != nil {
				n.logError("Error notifying that door 1 is open.", err.Error())
			}
			n.WasDoor1Open = true
		}
	}

	// Check how long Door2 has been open
	if room.Door2Closed {
		if n.WasDoor2Open {
			if err := n.sendMessage(fmt.Sprintf("%s's door is now closed.", room.Door2Name)); err != nil {
				n.logError("Error notifying that door 2 is now closed.", err.Error())
			}
			n.WasDoor2Open = false
		}
	} else {
		dur := time.Since(room.Door2StatusTime)
		if dur.Minutes() >= float64(config.DoorAlarmPeriod) {
			if err := n.sendMessage(fmt.Sprintf("%s's door has been open for %d minutes.", room.Door2Name, int(dur.Minutes()))); err != nil {
				n.logError("Error notifying that door 2 is open.", err.Error())
			}
			n.WasDoor2Open = true
		}
	}
}

// sendMessage sends the specified message to telegram
func (n *NotifyService) sendMessage(m string) error {
	c := telegram.Client{}
	return c.SendMessage(m)
}

// logError logs an error message to the logger
func (n *NotifyService) logError(v ...interface{}) {
	a := fmt.Sprint(v...)
	logger.Error("Notify [Err] ", a)
}
