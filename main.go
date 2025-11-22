package main

import (
	"prestasi_api/database"
)

func main() {
	database.ConnectPostgres()
	database.ConnectMongo()

	println("Prestasi API siap jalan...")
}
