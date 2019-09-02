package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
)

// Room holds the information about the room being monitored
type Room struct {
	Door1Name   string  `json:"door1name"` // Name of garage door 1
	Door1Open   bool    `json:"door1open"` // Whether door 1 is open
	Door2Name   string  `json:"door2name"` // Name of garage door 2
	Door2Open   bool    `json:"door2open"` // Whether door 2 is open
	Temperature float64 `json:"temp"`      // Room temperature
}

// OpenDoor issues the command to open the specified door number
func (r *Room) OpenDoor(doorNo int) error {
	if _, err := os.Stat("relay.py"); err != nil {
		return err
	}

	out, err := exec.Command("python", "relay.py", strconv.Itoa(doorNo)).CombinedOutput()
	if err != nil {
		msg := string(out)
		return errors.New(msg)
	}
	return err
}

// UpdateDoorStatus will update the Room telemetry with the new door statuses
func (r *Room) UpdateDoorStatus() error {
	r.logInfo("Updating door status")
	wd, err := os.Getwd()
	if err != nil {
		r.logError("Error getting current working directory.", err.Error())
		return err
	}
	dp := path.Join(wd, "data")

	if d1, err := r.readFileContents(path.Join(dp, "door1.state")); err != nil {
		r.logError("Failed to read door1 state.", err)
	} else {
		r.Door1Open = (d1 == "open")
	}
	if d2, err := r.readFileContents(path.Join(dp, "door2.state")); err != nil {
		r.logError("Failed to read door2 state.", err)
	} else {
		r.Door2Open = (d2 == "open")
	}
	return nil
}

func (r *Room) readFileContents(filePath string) (dir string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// logInfo logs an information message to the logger
func (r *Room) logInfo(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("Room: [Inf] ", a[1:len(a)-1])
}

// logError logs an error message to the logger
func (r *Room) logError(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Error("Room [Err] ", a[1:len(a)-1])
}
