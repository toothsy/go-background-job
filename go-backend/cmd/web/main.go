package main

import (
	"fmt"
	"github/toothsy/go-background-job/internal/config"
	"io"
	"log"
	"net/http"
	"os"
)

var app config.AppConfig
var portNumber = ":8080"

func main() {
	fmt.Println("connected")
	runner()
	// router := mux.NewRouter()
	fmt.Printf("Staring application on http://localhost%s", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func runner() {
	app.InProduction = false
	if app.InProduction {
		app.InfoLogger = log.New(io.Discard, "", 0)
		app.ErrorLogger = log.New(io.Discard, "", 0)
	} else {
		app.InfoLogger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		app.ErrorLogger = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)
	}

}
