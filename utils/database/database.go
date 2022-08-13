package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func CreateMongoClient(mongodb_url string) {
	log.Print("mongodb_url => ", mongodb_url)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodb_url))
	if err != nil {
		log.Panic(err)
	}
	Client = client
}
func GetMongoClient() *mongo.Client {
	return Client
}
