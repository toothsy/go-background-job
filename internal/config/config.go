package config

import (
	"github/toothsy/go-background-job/internal/workerpool"

	"go.mongodb.org/mongo-driver/mongo"
)

type AppConfig struct {
	InProduction  bool
	MongoDatabase *mongo.Database
	MongoClient   *mongo.Client
	WorkerPool    *workerpool.WorkerPool
}
