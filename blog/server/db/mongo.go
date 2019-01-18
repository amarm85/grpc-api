package db

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"os"
	"time"
)

func ConnectToMongoDB() *mongo.Client {

	mongoDBID := os.Getenv("MONGO_DB_ID")

	mongoDBPassword := os.Getenv("MONGO_DB_PASSWORD")

	if mongoDBPassword == "" || mongoDBID == "" {
		panic("Env variable MONGO_DB_PASSWORD or MONGO_DB_ID is not set")
	}
	url := fmt.Sprintf("mongodb://%s:%s@ds121222.mlab.com:21222/confusion", mongoDBID, mongoDBPassword)

	mClient, err := mongo.NewClient(url)

	if err != nil {
		log.Fatalf("Error in loading MangoDB driver %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	//defer mClient.Disconnect(ctx)
	err = mClient.Connect(ctx)

	if err != nil {
		log.Fatalf("Error in connecting MangoDB %v", err)
	}
	log.Printf("Successfully connected to MongoDB")
	return mClient
}
