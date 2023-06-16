package config

import (
	"github/toothsy/go-background-job/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type AppConfig struct {
	InProduction bool
	MonogoClient *mongo.Client
	WorkerPool   *models.WorkerPool
}
