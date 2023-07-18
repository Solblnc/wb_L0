package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Config struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
	SSLMode  string
}

type DataBase struct {
	client *sqlx.DB
}

func NewDataBase(cfg Config) (*DataBase, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode,
	)
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return &DataBase{}, fmt.Errorf("could not connect to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Query(`
    CREATE TABLE IF NOT EXISTS cars(
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL UNIQUE,
        url TEXT NOT NULL,
        rank INT NOT NULL
    );
    `)

	fmt.Println("created")

	res, err := db.Query(`
    SELECT * FROM cars;
    ;
    `)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

	return &DataBase{client: db}, nil
}
