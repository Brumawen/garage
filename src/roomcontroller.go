package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// RoomController handles the Web Methods for the room being monitored
type RoomController struct {
	Srv *Server
}

// AddController adds the controller routes to the router
func (c *RoomController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Methods("GET").Path("/room/get").Name("GetTelemetry").
		Handler(Logger(c, http.HandlerFunc(c.handleGetTelemetry)))
	router.Methods("POST").Path("/room/update").Name("UpdateTelemetry").
		Handler(Logger(c, http.HandlerFunc(c.handleUpdate)))
	router.Methods("POST").Path("/room/open/{doorNo}").Name("OpenDoor").
		Handler(Logger(c, http.HandlerFunc(c.handleOpenDoor)))
}

// handlerGetTelemetry will return the current telemetry for the room
func (c *RoomController) handleGetTelemetry(w http.ResponseWriter, r *http.Request) {
	if err := c.Srv.RoomService.UpdateTelemetry(); err != nil {
		http.Error(w, "Error updating telemetry", http.StatusInternalServerError)
	} else {
		if err := c.Srv.Room.WriteTo(w); err != nil {
			c.LogError("Error serializing telemetry.", err.Error)
			http.Error(w, "Error serializing telemetry", http.StatusInternalServerError)
		}
	}
}

// handleUpdate is called from the python script monitoring the door switches.  This call tells
// the server that the door status has changed
func (c *RoomController) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if err := c.Srv.RoomService.UpdateDoorStatus(); err != nil {
		http.Error(w, "Failed to update door status", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (c *RoomController) handleOpenDoor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	d := vars["doorNo"]
	doorNo := 0
	if d != "" {
		if i, err := strconv.Atoi(d); err == nil {
			doorNo = i
		}
	}
	if doorNo < 1 || doorNo > 2 {
		c.LogError("Invalid door number")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := c.Srv.RoomService.OpenDoor(doorNo)
	if err != nil {
		http.Error(w, "Failed", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

// LogInfo is used to log information messages for this controller.
func (c *RoomController) LogInfo(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("RoomController: [Inf] ", a[1:len(a)-1])
}

// LogError is used to log information messages for this controller.
func (c *RoomController) LogError(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("RoomController: [Err] ", a[1:len(a)-1])
}
