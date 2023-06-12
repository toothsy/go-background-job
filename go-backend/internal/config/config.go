package config

import (
	"github/toothsy/go-background-job/internal/models"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type AppConfig struct {
	InProduction bool
	InfoLogger   *log.Logger
	ErrorLogger  *log.Logger
	MonogoClient *mongo.Client
	WorkerPool   *models.WorkerPool
}
