package db

import (
	"context"
	"fmt"
	"hub-api/config"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func InitDB() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.DatabaseURL))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected to db succesfully")

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Ping to db succesful")

	// get the database
	DB = client.Database("pixelhub")

}
