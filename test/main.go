package main

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://192.168.8.133:1883")
	opts.SetUsername("mqttuser")
	opts.SetPassword("1qazxsw@")

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Print("publishing")
	token := client.Publish("home/garage/temperature", byte(0), false, 24.25)
	token.Wait()

	client.Disconnect(250)
}
