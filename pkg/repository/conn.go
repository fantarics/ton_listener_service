package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func NewDB(connection string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: connection,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
