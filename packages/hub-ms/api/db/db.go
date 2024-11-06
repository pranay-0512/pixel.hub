package db

import (
	"context"
	"fmt"
	"hub-api/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func InitDB(ctx context.Context) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.DatabaseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("pinging failed: %w", err)
	}
	fmt.Println("Connection to mongodb successful")
	DB = client.Database(config.AppConfig.DatabaseName)

	return DB, nil
}

func DisconnectDB(ctx context.Context, client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect mongoDB: %w", err)
	}
	fmt.Println("Disconnected from mongodb")
	return nil
}
