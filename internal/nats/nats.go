package nats

import (
	"L0/internal/database"
	"L0/internal/model"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type Nats interface {
	LoadCacheFromDB() error
	GetOrderFromCache(id int) (model.Order, error)
}

type NatsServer struct {
	DB    database.Repository
	Cache map[int]model.Order
}

func NewNatsServer(DB database.Repository, Cache map[int]model.Order) *NatsServer {
	return &NatsServer{DB: DB, Cache: Cache}
}

func (ns *NatsServer) NatsConnect() error {

	if err := ns.LoadCacheFromDB(); err != nil {
		return fmt.Errorf("cannot load cache from db: %w", err)
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return fmt.Errorf("cannot connect to nats-server: %w", err)
	}

	js, _ := nc.JetStream()
	js.AddStream(&nats.StreamConfig{Name: "wb", Subjects: []string{"order"}})
	str, _ := js.StreamInfo("wb")
	if str == nil {
		fmt.Println("no stream")
	}

	sub, err := js.SubscribeSync("order")
	if err != nil {
		return fmt.Errorf("cannot create subscription: %w", err)
	}

	log.Println("connected to stream")

	for {
		msg, err := sub.NextMsg(1 * time.Second)
		if err != nil {
			continue
		}

		if err = ns.GetMessage(msg); err != nil {
			log.Println(err)
		}

	}

	return nil

}

func (ns *NatsServer) LoadCacheFromDB() error {
	cache, err := ns.DB.LoadToCache()
	if err != nil {
		return err
	}

	ns.Cache = cache
	log.Println("cache loaded from database")

	return nil
}

func (ns *NatsServer) GetOrderFromCache(id int) (model.Order, error) {
	order, ok := ns.Cache[id]
	if !ok {
		return model.Order{}, fmt.Errorf("cannot find order with given id")
	}

	return order, nil
}

func (ns *NatsServer) GetMessage(msg *nats.Msg) error {
	var order model.Order

	if err := json.Unmarshal(msg.Data, &order); err != nil {
		return fmt.Errorf("there is no order data in message: %s", string(msg.Data))
	}

	ID, err := ns.DB.Save(order)
	if err != nil {
		return fmt.Errorf("cannot save order in db: %w", err)
	}

	ns.Cache[ID] = order

	return nil
}
