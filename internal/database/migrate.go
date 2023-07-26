package database

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func (d *DataBase) Migrate() error {

	driver, err := postgres.WithInstance(d.client.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("can not create a postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "orders", driver)
	if err != nil {
		return fmt.Errorf("file error: %w", err)
	}

	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {

			return fmt.Errorf("could not run up migrations: %w", err)
		}
	}

	log.Println("Successfully migrated our database")

	return nil

}

func (d *DataBase) MigrateDown() error {
	fmt.Println("migrating down  our database")

	driver, err := postgres.WithInstance(d.client.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("can not create a postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "wb", driver)
	if err != nil {
		return fmt.Errorf("file error: %w", err)
	}

	if err = m.Down(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {

			return fmt.Errorf("could not run up migrations: %w", err)
		}
	}

	fmt.Println("Successfully migrated DOWN our database")

	return nil

}
