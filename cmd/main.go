package main

import (
	"L0/internal/database"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
)

func main() {
	//nc, err := nats.Connect(nats.DefaultURL)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//js, _ := nc.JetStream()
	//js.AddStream(&nats.StreamConfig{Name: "test", Subjects: []string{"bye"}})
	//str, _ := js.StreamInfo("test")
	//if str == nil {
	//	fmt.Println("no stream")
	//}
	//
	//js.Publish("bye", []byte("bye from Emil upgraded"))
	//js.Subscribe("bye", func(msg *nats.Msg) {
	//	msg.Ack()
	//	fmt.Println(string(msg.Data))
	//})

	db, err := database.NewDataBase(database.Config{
		Host:     "localhost",
		Port:     "5433",
		User:     "postgres",
		DBName:   "test",
		Password: "postgres",
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db)
	fmt.Println("hello")

}
