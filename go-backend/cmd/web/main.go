package main

import (
	"context"
	"fmt"
	"github/toothsy/go-background-job/internal/config"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
	app.MongoURI = "mongodb+srv://user:ctYORHp3WxOqijzw@background-job-go.kjbmtqg.mongodb.net/"
	client, err := mongo.NewClient(options.Client().ApplyURI(app.MongoURI))
	if err != nil {
		app.ErrorLogger.Fatal(err)
		return nil, err
	}

	return client, nil
}
