package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	SecretsCollection *mongo.Collection
	Ctx               = context.TODO()
)

func Connect() {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		log.Fatal(err)
	}

	Ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(Ctx)
	if err != nil {
		log.Fatal(err)
	}

	secretDatabase := client.Database(os.Getenv("MONGO_DATABASE"))
	SecretsCollection = secretDatabase.Collection(os.Getenv("MONGO_COLLECTION"))
}
