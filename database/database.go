package database

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var once sync.Once

// Database - Database Singleton that holds database reference
type Database struct {
	client *mongo.Client
}

var instance *Database

// GetInstance - returns instance of database
func GetInstance() *Database {
	once.Do(func() {
		instance = &Database{}

		clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI_DEV"))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var err error
		instance.client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err)
		}
		defer cancel()
	})

	return instance
}
