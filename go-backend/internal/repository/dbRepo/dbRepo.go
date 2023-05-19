package dbrepo

import (
	"github/toothsy/go-background-job/internal/config"
	"github/toothsy/go-background-job/internal/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDBRepo struct {
	AppConfig *config.AppConfig
	DB        *mongo.Database
}

func NewMongoConnection(a *config.AppConfig, m *mongo.Database) repository.DatabaseRepo {
	return &mongoDBRepo{
		AppConfig: a,
		DB:        m,
	}

}
