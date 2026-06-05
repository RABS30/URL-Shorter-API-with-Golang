package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func DatabaseConnect() *pgxpool.Pool {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("cannot load .env, " + err.Error())
	}

	db, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		panic("cannot connect to database" + err.Error())
	}

	return db
}
