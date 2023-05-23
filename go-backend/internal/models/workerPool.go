package models

import (
	"github/toothsy/go-background-job/internal/config"
	"sync/atomic"
)

type WorkerPool struct {
	jobCh  chan Job
	config *config.WorkerPoolConfig
	done   atomic.Bool
}

// Enqueue adds the gives job to pool
func (wp *WorkerPool) Enqueue(job Job) {
	if wp.jobCh != nil {
		wp.jobCh <- job
	}
}

// Shutdown closes the jobs channel
func (wp *WorkerPool) Shutdown() {
	close(wp.jobCh)
	wp.done.Store(true)
}

// delegates the work to worker routines
func (wp *WorkerPool) runWorker() {
	for {
		// Continuously check for new jobs till shutdown signal
		select {
		case job := <-wp.jobCh:
			wp.ProcessJob(job)
		default:
			if wp.jobCh == nil || wp.done.Load() {
				return
			}
		}
	}
}

// Run spawns go routines for the worker pool to handle jobs
func (wp *WorkerPool) Run() {
	for i := 0; i < wp.config.MaxWorkers; i++ {
		go wp.runWorker()
	}
}

// ProcessJob puts the image from job to the database
func (wp *WorkerPool) ProcessJob(dequedJob Job) {
	// two kinds of jobs one to insert, one to lookup the user credential
}
