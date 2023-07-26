package main

import (
	"L0/internal/model"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func main() {

	order, err := model.NewOrder("model.json")
	if err != nil {
		log.Fatalf("cannot read a model, %s", err)
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("cannot marshall order, %s", err)
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("cannot connect to nats-server: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	msg := &nats.Msg{
		Subject: "order",
		Data:    data,
	}

	if _, err := js.PublishMsg(msg); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Message published")

	//nc.Publish("order", data)
	//time.Sleep(1 * time.Second)

}
