package main

import (
	"github/toothsy/go-background-job/internal/config"
	"github/toothsy/go-background-job/internal/driver"
	"github/toothsy/go-background-job/internal/handlers"
	"github/toothsy/go-background-job/internal/models"
	dbrepo "github/toothsy/go-background-job/internal/repository/dbRepo"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
)

var app config.AppConfig
var portNumber = ":8080"
var mongoUri string

func main() {
	cancel, err := runner()
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Disconnect(&app)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Staring application on http://localhost%s", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func runner() (func(), error) {
	app.InProduction = false

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
		log.Fatal(err)
		return nil, err
	}
	mongoUri = os.Getenv("MONGO_URI")
	cancel, err := driver.ConnectMongoDB(&app, mongoUri)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// allowing handlers the access to appconfig

	repo := handlers.NewRepo(&app)
	// initiating a new repo instance, so that handler functions can use app config
	handlers.NewHandlers(repo)

	// allowing dbRepo the access to appconfig
	dbrepo.NewMongoConnection(&app)
	return cancel, nil
}
