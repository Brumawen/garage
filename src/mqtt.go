package main

import (
	"errors"
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Mqtt publishes the telemetry to a MQTT Broker and
// subscribes to commands
type Mqtt struct {
	Srv               *Server     // Server instance
	LastUpdateAttempt time.Time   // Last time an update was attempted
	LastUpdate        time.Time   // Last time an update was published
	client            MQTT.Client // MQTT client
	ignoreCommands    bool        // Signals that commands must be ignored
}

// Initialize initializes the MQTT client
func (m *Mqtt) Initialize() error {
	if !m.Srv.Config.EnableMqtt {
		m.logInfo("MQTT has been disabled")
		return nil
	}
	if m.Srv.Config.MqttHost == "" {
		m.logError("MQTT Host has not been configured.")
		m.Srv.Config.EnableMqtt = false
		return errors.New("host has not been configured")
	}
	if m.Srv.Config.MqttUsername == "" {
		m.logError("MQTT Username has not been configured.")
		m.Srv.Config.EnableMqtt = false
		return errors.New("username has not been configured")
	}
	if m.Srv.Config.MqttPassword == "" {
		m.logError("MQTT Password has not been configured.")
		m.Srv.Config.EnableMqtt = false
		return errors.New("password has not been configured")
	}

	// Connect and send meta information
	m.logInfo("Connecting to the MQTT Broker.")
	m.ignoreCommands = true

	opts := MQTT.NewClientOptions()
	opts.AddBroker(m.Srv.Config.MqttHost)
	opts.SetUsername(m.Srv.Config.MqttUsername)
	opts.SetPassword(m.Srv.Config.MqttPassword)

	opts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		m.logError("Disconnected from MQTT Broker.", err.Error())
	})
	opts.SetOnConnectHandler(func(client MQTT.Client) {
		m.logInfo("Connected to the MQTT Broker. Subscribing to topics.")
		if token := client.Subscribe("home/garage/door1/set", byte(1), nil); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		if token := client.Subscribe("home/garage/door2/set", byte(1), nil); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		m.logInfo("Subscription complete.")
	})
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		m.logInfo("Command received.", msg.Topic(), string(msg.Payload()))
		if m.ignoreCommands {
			m.logInfo("Commands are currently being ignored")
			return
		}
		if msg.Topic() == "home/garage/door1/set" {
			pl := string(msg.Payload())
			m.logInfo("Received Door 1 Set command with payload of:", pl)
			if pl == "ON" {
				// Check if the door is open and close it
				if !m.Srv.Room.Door1Closed {
					m.logInfo("Closing door 1")
					m.Srv.RoomService.OpenDoor(1)
				}
			} else if pl == "OFF" {
				// Check if the door is closed an open it
				if m.Srv.Room.Door1Closed {
					m.logInfo("Opening door 1")
					m.Srv.RoomService.OpenDoor(1)
				}
			}
		} else if msg.Topic() == "home/garage/door2/set" {
			pl := string(msg.Payload())
			m.logInfo("Received Door 2 Set command with payload of:", pl)
			if pl == "ON" {
				// Check if the door is open and close it
				if !m.Srv.Room.Door2Closed {
					m.logInfo("Closing door 2")
					m.Srv.RoomService.OpenDoor(2)
				}
			} else if pl == "OFF" {
				// Check if the door is closed an open it
				if m.Srv.Room.Door2Closed {
					m.logInfo("Opening door 2")
					m.Srv.RoomService.OpenDoor(2)
				}
			}
		}
	})

	m.client = MQTT.NewClient(opts)
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		m.logError("Error connecting to MQTT Broker.", token.Error())
		return token.Error()
	}

	return nil
}

// Close closes the MQTT client and disconnects
func (m *Mqtt) Close() {
	m.client.Disconnect(250)
}

// SendTelemetry sends the current states of the devices to the MQTT Broker
func (m *Mqtt) SendTelemetry() error {
	if !m.Srv.Config.EnableMqtt {
		return nil
	}

	m.logInfo("Publishing telemetry to MQTT")
	m.LastUpdateAttempt = time.Now()

	// DOOR 1
	doorState := "OFF"
	if m.Srv.Room.Door1Closed {
		doorState = "ON"
	}
	m.logInfo("Publishing door1 state")
	token := m.client.Publish("home/garage/door1", byte(0), true, doorState)
	if token.Wait() && token.Error() != nil {
		m.logError("Error sending door 1 state to MQTT Broker.", token.Error())
		return token.Error()
	}

	// DOOR 2
	doorState = "OFF"
	if m.Srv.Room.Door2Closed {
		doorState = "ON"
	}
	m.logInfo("Publishing door2 state")
	token = m.client.Publish("home/garage/door2", byte(0), true, doorState)
	if token.Wait() && token.Error() != nil {
		m.logError("Error sending door 2 state to MQTT Broker.", token.Error())
		return token.Error()
	}

	// Temperature
	m.logInfo("Publishing temperature")
	token = m.client.Publish("home/garage/temperature", byte(0), true, fmt.Sprintf("%.1f", m.Srv.Room.Temperature))
	if token.Wait() && token.Error() != nil {
		m.logError("Error sending temperature state to MQTT Broker.", token.Error())
		return token.Error()
	}

	m.LastUpdate = time.Now()
	m.ignoreCommands = false

	return nil
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
