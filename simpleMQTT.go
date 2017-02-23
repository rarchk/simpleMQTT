package main 

import(
	"fmt"
	// import from eclipse Paho go MQTT lib
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
  "net/url"
  
)

func main() {
    go func(topic string) {
        opts := createClientOptions("sub", os.Getenv("CLOUDMQTT_URL"));
        fmt.Println(opts);
        client := MQTT.NewClient(opts)
        client.Start()


        t, _ := MQTT.NewTopicFilter(topic, 0)
        client.StartSubscription(func(client *MQTT.MQTTClient, msg MQTT.Message) {
            fmt.Println("Topic=", msg.Topic(), "Payload=", string(msg.Payload()))
        }, t)
    }("#")

    timer := time.NewTicker(1 * time.Second)
    opts := createClientOptions("pub", os.Getenv("CLOUDMQTT_URL"))
    client := MQTT.NewClient(opts)
    client.Start()

    for t := range timer.C {
        client.Publish(0, "currentTime", t.String())
    }
}

func createClientOptions(clientId, raw string) *MQTT.ClientOptions {
    uri, _ := url.Parse(raw)
    opts := MQTT.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
    opts.SetUsername(uri.User.Username())
    password, _ := uri.User.Password()
    opts.SetPassword(password)
    opts.SetClientId(clientId)

    return opts
}

