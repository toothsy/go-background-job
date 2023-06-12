package models

import (
	"github/toothsy/go-background-job/internal/constants"
	"log"
	"sync/atomic"
	"time"
)

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

type WorkerPool struct {
	JobCh  chan Job
	Config *WorkerPoolConfig
	Done   *atomic.Bool
}

// Enqueue adds the gives job to pool
func (wp *WorkerPool) Enqueue(job Job) {
	if wp.JobCh != nil {
		wp.JobCh <- job
	}
}

// Shutdown closes the jobs channel
func (wp *WorkerPool) Shutdown() {
	close(wp.JobCh)
	wp.Done.Store(true)
}

// delegates the work to worker routines
func (wp *WorkerPool) runWorker() {
	for {
		// Continuously check for new jobs till shutdown signal
		select {
		case job := <-wp.JobCh:
			wp.ProcessJob(job)
		default:
			if wp.JobCh == nil || wp.Done.Load() {
				log.Println("waiting on job")
				return
			}
		}
	}
}

// Run spawns go routines for the worker pool to handle jobs
func (wp *WorkerPool) Run() {
	for i := 0; i < wp.Config.MaxWorkers; i++ {
		go wp.runWorker()

	}
}

// ProcessJob puts the image from job to the database
func (wp *WorkerPool) ProcessJob(dequedJob Job) {
	// two kinds of jobs one to insert, one to lookup the user credential
	if dequedJob.JobType == constants.Authenticate {
		handleAuth(dequedJob)
	} else if dequedJob.JobType == constants.Upload {
		handleUpload(dequedJob)
	}
}
func handleAuth(dequedJob Job) {
	log.Println("got the auth job", dequedJob)

}
func handleUpload(dequedJob Job) {
	log.Println("got the upload job", dequedJob)
}
