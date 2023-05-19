package models

import "github/toothsy/go-background-job/internal/repository"

// constants for job status
const (
	Queued = iota
	Running
	Completed
	Failed
)

//QueueManager holds the DbRepo,workerPool and JobQueue,QueueConfig
type QueueManager struct {
	DatabaseRepo repository.DatabaseRepo
}

//Job defines what workers pass around to execute and load to db
type Job struct {
	Id     string
	Status Status
	Image  []byte
}
