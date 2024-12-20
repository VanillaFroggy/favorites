package db

import (
	"favorites/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func ConnectDB() (*sqlx.DB, error) {
	cfg := config.LoadConfig()
	db, err := sqlx.Connect("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal("Failed to connect to DB: " + err.Error())
		return nil, err
	}
	return db, err
}
