package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pg *pgxpool.Pool

func ConnectPostgres() error {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, dbname,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("gagal konek Postgres: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("Postgres tidak merespons: %v", err)
	}

	Pg = pool
	log.Println("✅ PostgreSQL terhubung")
	return nil
}
