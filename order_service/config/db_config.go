package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DBClient      *mongo.Client
	MONGO_DB_NAME = "pizza-shop-eda"
)

func init() {
	InitializeDB()
}

func InitializeDB() (*mongo.Client, error) {
	var err error
	logger.Log("Initializin database once more")
	if DBClient == nil {
		DBClient, err = initDatabase()
		if err != nil {
			log.Fatal("Failde to initialize database: %v", err)
		}
	}
	return DBClient, err
}

func initDatabase() (*mongo.Client, error) {
	dbURL := GetEnvProperty("database_url")
	if dbURL == "" {
		return nil, fmt.Errorf("database_url is not set in the environment variable")
	}
	clientOptions := options.Client().ApplyURI(dbURL).
		SetMaxPoolSize(600).
		SetMinPoolSize(50).
		SetMaxConnIdleTime(30 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize database: %v", err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("MongoDB is unreachable: %v", err)
	}
	log.Printf("Connected to MongoDB: %s", dbURL)
	return client, nil
}
