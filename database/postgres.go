package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pg *pgxpool.Pool

func ConnectPostgres() {
	dsn := "postgres://postgres:12345678@localhost:5432/prestasi_api_db?sslmode=disable"


	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Gagal konek Postgres:", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Postgres tidak merespons:", err)
	}

	Pg = pool
	fmt.Println("PostgreSQL terhubung")
}
