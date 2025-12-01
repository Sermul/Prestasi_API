package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo *mongo.Database

func ConnectMongo() error {
	uri := os.Getenv("MONGO_URI")
	dbname := os.Getenv("MONGO_DB")

	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		return fmt.Errorf("gagal konek MongoDB: %v", err)
	}

	Mongo = client.Database(dbname)
	log.Println("✅ MongoDB terhubung")
	return nil
}
