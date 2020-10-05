package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetConnection() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://SwarzChen:aa889837088@cluster0-shard-00-00.noe6r.gcp.mongodb.net:27017,cluster0-shard-00-01.noe6r.gcp.mongodb.net:27017,cluster0-shard-00-02.noe6r.gcp.mongodb.net:27017/<dbname>?ssl=true&replicaSet=atlas-14mf0y-shard-0&authSource=admin&retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	return client
}
