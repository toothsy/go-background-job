package config

import (
	"log"
	"time"
)

type AppConfig struct {
	InProduction bool
	InfoLogger   *log.Logger
	ErrorLogger  *log.Logger
	MongoURI     string
}

type WorkerPoolConfig struct {
	MaxWorkers     int
	MaxQueueSize   int
	MaxRetries     int
	RetryDelay     time.Duration
	Timeout        time.Duration
	MetricsTricker time.Duration
	PruneInterval  time.Duration
	PanicHandler   func(job interface{})
}
