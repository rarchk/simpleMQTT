package main 

import(
	"fmt"
	// import from eclipse Paho go mqtt lib
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)


//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
  fmt.Printf("TOPIC: %s\n", msg.Topic())
  fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
  //create a ClientOptions struct setting the broker address, clientid, turn
  //off trace output and set the default message handler
  opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
  opts.SetClientID("go-simple")
  opts.SetDefaultPublishHandler(f)

  //create and start a client using the above ClientOptions
  c := MQTT.NewClient(opts)
  if token := c.Connect(); token.Wait() && token.Error() != nil {
    panic(token.Error())
  }

  //subscribe to the topic /go-mqtt/sample and request messages to be delivered
  //at a maximum qos of zero, wait for the receipt to confirm the subscription
  if token := c.Subscribe("room1/readings/temperature", 0, nil); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    os.Exit(1)
  }

  //Publish 5 messages to room1/readings/temperature at qos 1 and wait for the receipt
  //from the server after sending each message
  for i := 0; i < 5; i++ {
    text := fmt.Sprintf("this is msg #%d!", i)
    token := c.Publish("room1/readings/temperature", 0, false, text)
    token.Wait()
  }

  time.Sleep(3 * time.Second)

  //unsubscribe from room1/readings/temperature
  if token := c.Unsubscribe("room1/readings/temperature"); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    os.Exit(1)
  }

  c.Disconnect(250)
}
