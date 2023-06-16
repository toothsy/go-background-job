package dbrepo

import (
	"github/toothsy/go-background-job/internal/config"
	"github/toothsy/go-background-job/internal/repository"
)

type mongoDBRepo struct {
	AppConfig *config.AppConfig
}

func NewMongoConnection(a *config.AppConfig) repository.DatabaseRepo {
	return &mongoDBRepo{
		AppConfig: a,
	}
}
