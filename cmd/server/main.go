package main

import (
	"L0/internal/config"
	"L0/internal/database"
	"L0/internal/model"
	"L0/internal/nats"
	"L0/internal/server"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
)

func main() {

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDataBase(database.Config{
		Host:     cfg.Config.Host,
		Port:     cfg.Config.Port,
		User:     cfg.Config.User,
		DBName:   cfg.Config.DBName,
		Password: cfg.Config.Password,
		SSLMode:  cfg.Config.SSLMode,
	})

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Migrate(); err != nil {
		log.Fatalf("failed to migrate: %w", err)
	}

	//cache, err := db.LoadToCache()
	//if err != nil {
	//	log.Fatal(err)
	//}

	cache := make(map[int]model.Order)

	ns := nats.NewNatsServer(db, cache)

	server := server.NewHandler(ns)

	r := server.Router()

	restServer := &http.Server{Addr: cfg.Config.RestPort, Handler: r}

	go restServer.ListenAndServe()

	log.Println("Server is running on port :8080")

	if err = ns.NatsConnect(); err != nil {
		log.Fatal(err)
	}

}
