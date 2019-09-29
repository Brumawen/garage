package main

import (
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Mqtt publishes the telemetry to a MQTT Broker and
// subscribes to commands
type Mqtt struct {
	Srv        *Server
	LastUpdate time.Time
	IsRunning  bool
}

// Initialize initializes the MQTT client
func (m *Mqtt) Initialize() {
	if !m.Srv.Config.EnableMqtt {
		m.logInfo("MQTT has been disabled")
		return
	}
	if m.Srv.Config.MqttHost == "" {
		m.logError("MQTT Host has not been configured.")
		m.Srv.Config.EnableMqtt = false
		return
	}
	if m.Srv.Config.MqttUsername == "" {
		m.logError("MQTT Username has not been configured.")
		m.Srv.Config.EnableMqtt = false
		return
	}
	if m.Srv.Config.MqttPassword == "" {
		m.logError("MQTT Password has not been configured.")
		m.Srv.Config.EnableMqtt = false
		return
	}

	// Connect and send meta information

}

// SendTelemetry sends the current states of the devices to the MQTT Broker
func (m *Mqtt) SendTelemetry() error {
	if !m.Srv.Config.EnableMqtt {
		return nil
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker(m.Srv.Config.MqttHost)
	opts.SetUsername(m.Srv.Config.MqttUsername)
	opts.SetPassword(m.Srv.Config.MqttPassword)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		m.logError("Error connecting to MQTT Broker.", token.Error())
		return token.Error()
	}

	m.logInfo("Publishing telemetry to MQTT")

	// DOOR 1
	doorState := "OFF"
	if m.Srv.Room.Door1Closed {
		doorState = "ON"
	}
	token := client.Publish("home/garage/door1", byte(0), true, doorState)
	if token.Wait() && token.Error() != nil {
		m.logError("Error sending door 1 state to MQTT Broker.", token.Error())
		client.Disconnect(250)
		return token.Error()
	}

	// DOOR 2
	doorState = "OFF"
	if m.Srv.Room.Door2Closed {
		doorState = "ON"
	}
	token = client.Publish("home/garage/door2", byte(0), true, doorState)
	if token.Wait() && token.Error() != nil {
		m.logError("Error sending door 2 state to MQTT Broker.", token.Error())
		client.Disconnect(250)
		return token.Error()
	}

	// Temperature
	token = client.Publish("home/garage/temperature", byte(0), true, m.Srv.Room.Temperature)
	if token.Wait() && token.Error() != nil {
		m.logError("Error sending temperature state to MQTT Broker.", token.Error())
		client.Disconnect(250)
		return token.Error()
	}

	client.Disconnect(250)
	return nil
}

// subscribe listens for commands from the MQTT broker
func (m *Mqtt) subscribe() {

}

// logInfo logs an information message to the logger
func (m *Mqtt) logInfo(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("Mqtt: [Inf] ", a[1:len(a)-1])
}

// logError logs an error message to the logger
func (m *Mqtt) logError(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Error("Mqtt [Err] ", a[1:len(a)-1])
}
