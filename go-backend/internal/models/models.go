package models

import (
	"fmt"
	"github/toothsy/go-background-job/internal/repository"
	"mime/multipart"
)

// QueueManager holds the DbRepo,workerPool and JobQueue,QueueConfig
type QueueManager struct {
	DatabaseRepo repository.DatabaseRepo
	WorkerPool   WorkerPool
}

// Job defines what workers pass around to execute and load to db
type Job struct {
	Id           string
	Status       int // refer to contants
	JobType      int // refer to contants
	Image        *multipart.FileHeader
	Username     string
	PasswordHash string
}

func (j Job) String() string {
	return fmt.Sprintf(`
	ID:%s
	Status:%d
	JobType:%d
	ImageNil:%t
	Username:%s
	PassWordHash:%s`, j.Id, j.Status, j.JobType, j.Image == nil, j.Username, j.PasswordHash)
}
