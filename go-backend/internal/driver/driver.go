package driver

import (
	"context"
	"github/toothsy/go-background-job/internal/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connects to the mongo db instance, returns cancel function and error
func ConnectMongoDB(app *config.AppConfig, mongoURI string) (func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		cancel()
		return func() {}, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		cancel()
		return func() {}, err

	}
	log.Println("Ping to Mongo DB successful")
	app.MonogoClient = client
	return cancel, nil
}

// Disconnect disconnects the mongo db client instance
func Disconnect(app *config.AppConfig) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	app.MonogoClient.Disconnect(ctx)

}
