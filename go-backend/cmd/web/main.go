package main

import (
	"context"
	"fmt"
	"github/toothsy/go-background-job/internal/config"
	"github/toothsy/go-background-job/internal/handlers"
	"github/toothsy/go-background-job/internal/models"
	"io"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var app config.AppConfig
var portNumber = ":8080"

func main() {
	mongoClient, err := runner()
	if err != nil {
		app.ErrorLogger.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer mongoClient.Disconnect(ctx)
	defer cancel()

	if err != nil {
		app.ErrorLogger.Fatal(err)
	}
	fmt.Printf("Staring application on http://localhost%s", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func runner() (*mongo.Client, error) {
	app.InProduction = false
	if app.InProduction {
		app.InfoLogger = log.New(io.Discard, "", 0)
		app.ErrorLogger = log.New(io.Discard, "", 0)
	} else {
		app.InfoLogger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		app.ErrorLogger = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)
	}

	WorkerConfig := &models.WorkerPoolConfig{
		MaxWorkers:   2,
		MaxQueueSize: 10,
		MaxRetries:   5,
		RetryDelay:   time.Millisecond * 500,
		Timeout:      time.Minute,
	}
	jobChannel := make(chan models.Job, WorkerConfig.MaxQueueSize)
	app.WorkerPool = &models.WorkerPool{
		JobCh:  jobChannel,
		Config: WorkerConfig,
		Done:   &atomic.Bool{},
	}
	app.WorkerPool.Run()
	err := godotenv.Load("./secret.env")
	if err != nil {
		app.ErrorLogger.Fatal(err)
		return nil, err
	}
	mongoUri := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		app.ErrorLogger.Fatal(err)
		return nil, err
	}
	app.MonogoClient = client
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	return client, nil
}
