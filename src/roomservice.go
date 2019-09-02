package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"

	gopitools "github.com/brumawen/gopi-tools/src"
)

// RoomService contains service methods for the room being monitored
type RoomService struct {
	Srv *Server
}

// OpenDoor issues the command to open the specified door number
func (r *RoomService) OpenDoor(doorNo int) error {
	if _, err := os.Stat("relay.py"); err != nil {
		r.logError("File relay.py does not exist.")
		return err
	}

	out, err := exec.Command("python", "relay.py", strconv.Itoa(doorNo)).CombinedOutput()
	if err != nil {
		r.logError("Failed to open door", doorNo, err.Error)
		msg := string(out)
		return errors.New(msg)
	}
	return err
}

// UpdateTelemetry will update all telemetry associated with the room
func (r *RoomService) UpdateTelemetry() error {
	// Get the temperature probe
	tmp := gopitools.OneWireTemp{}
	defer tmp.Close()
	tmp.ID = ""

	// Get the available one-wire devices
	r.logDebug("Getting one-wire device list.")
	devlst, err := gopitools.GetDeviceList()
	if err != nil {
		msg := "Error getting one-wire device list. " + err.Error() + "."
		r.logError(msg)
		return err
	} else {
		if len(devlst) == 0 {
			msg := "No temperature device found. Cable could be disconnected."
			r.logError(msg)
		} else {
			r.logDebug("Reading temperature from ", devlst[0].Name)
			tmp.ID = devlst[0].ID
			temp, err := tmp.ReadTemp()
			if err != nil {
				msg := "Error reading temperature. " + err.Error() + "."
				r.logError(msg)
				return err
			} else {
				r.Srv.Room.Temperature = temp
			}
		}
	}

	return nil
}

// UpdateDoorStatus will update the Room telemetry with the new door statuses
func (r *RoomService) UpdateDoorStatus() error {
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
		r.Srv.Room.Door1Open = (d1 == "open")
	}
	if d2, err := r.readFileContents(path.Join(dp, "door2.state")); err != nil {
		r.logError("Failed to read door2 state.", err)
	} else {
		r.Srv.Room.Door2Open = (d2 == "open")
	}
	return nil
}

func (r *RoomService) readFileContents(filePath string) (dir string, err error) {
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

// logDebug logs a debug message to the logger
func (r *RoomService) logDebug(v ...interface{}) {
	if r.Srv.VerboseLogging {
		a := fmt.Sprint(v)
		logger.Info("Server: [Dbg] ", a[1:len(a)-1])
	}
}

// logInfo logs an information message to the logger
func (r *RoomService) logInfo(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("RoomService: [Inf] ", a[1:len(a)-1])
}

// logError logs an error message to the logger
func (r *RoomService) logError(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Error("Room [Err] ", a[1:len(a)-1])
}
