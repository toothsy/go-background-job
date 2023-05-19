package models

import "github/toothsy/go-background-job/internal/config"

type WorkerPool struct {
	jobs   chan Job
	config *config.WorkerPoolConfig
}
