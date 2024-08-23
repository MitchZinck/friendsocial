package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	DB_USER     = "mitchell.zinck"
	DB_PASSWORD = ""
	DB_NAME     = "friendsocialdb"
)

var DB *pgxpool.Pool

func InitDB() {
	var err error
	psqlInfo := fmt.Sprintf("postgresql://%s:%s@localhost:5432/%s?sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	config, err := pgxpool.ParseConfig(psqlInfo)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	DB, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
