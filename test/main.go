package main

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	quit := make(chan string)
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://172.16.30.114:1883")
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		fmt.Println("Topic", msg.Topic(), "Payload", string(msg.Payload()))
		if string(msg.Payload()) == "quit" {
			quit <- "quit"
		}
	})
	opts.SetOnConnectHandler(func(client MQTT.Client) {
		fmt.Println("Connected")
		fmt.Println("subscribing")
		if token := client.Subscribe("home/test", byte(1), nil); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		if token := client.Subscribe("home/test1", byte(1), nil); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	})
	opts.SetConnectionLostHandler(func(client MQTT.Client, e error) {
		fmt.Println("Disconnected")
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	_ = <-quit

	client.Disconnect(250)
}
