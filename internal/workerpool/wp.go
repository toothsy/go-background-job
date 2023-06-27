package workerpool

import (
	"context"
	"crypto/subtle"
	"github/toothsy/go-background-job/internal/constants"
	"github/toothsy/go-background-job/internal/models"
	"github/toothsy/go-background-job/internal/repository"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
	JobCh  chan *models.Job
	Config *WorkerPoolConfig
	Done   *atomic.Bool
}

var JobContextMap *sync.Map
var DB repository.DatabaseRepo
var mongoDatbase *mongo.Database

func (wp *WorkerPool) Init(contextMap *sync.Map, repo repository.DatabaseRepo, mdb *mongo.Database) {
	JobContextMap = contextMap
	DB = repo
	mongoDatbase = mdb
}

// Enqueue adds the gives job to pool
func (wp *WorkerPool) Enqueue(job *models.Job) {
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
		// Continuously check for new job submission on channel till shutdown signal
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
		log.Println("fired routines")
	}
}

// ProcessJob puts the image from job to the database
func (wp *WorkerPool) ProcessJob(dequedJob *models.Job) {
	if dequedJob == nil {
		return
	}
	// two kinds of jobs one to insert, one to lookup the user credential
	log.Println("processing")
	if dequedJob.JobType == constants.Authenticate {
		handleAuth(dequedJob)
	} else if dequedJob.JobType == constants.Upload {
		handleUpload(dequedJob)
	}
}
func handleAuth(dequedJob *models.Job) {
	// log.Println("got the auth job", dequedJob)
	dbUser, err := DB.FetchUser(&dequedJob.User)
	if err != nil {
		log.Println("error in fetching user from authentication ", err)
	}
	updateClient(dequedJob.Id, constants.Running, gin.H{"message": "job started", "code": 200})
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(dequedJob.User.PasswordHash))
	isCorrectPassword := err == nil

	_, ok := JobContextMap.Load(dequedJob.Id)
	if !ok {
		log.Println("non-existent job retrieved from authentication ", err)

	}
	if !ok {
		log.Println("attempted to cast non-gin-context object from authentication ", err)

	}
	if isCorrectPassword {
		updateClient(dequedJob.Id, constants.Completed, gin.H{"message": "proper password"})
	} else {
		updateClient(dequedJob.Id, constants.Failed, gin.H{"message": "improper proper password"})
	}

}
func handleUpload(dequedJob *models.Job) {
	log.Println("got the upload job", dequedJob)
	cInterafce, ok := JobContextMap.Load(dequedJob.Id)
	if !ok {
		log.Println("error in retrieving gin context ")
	}
	c, ok := cInterafce.(*gin.Context)
	if !ok {
		log.Println("error in casting gin context ")
	}

	_, err := mongoDatbase.Collection("imageData").InsertOne(context.Background(), dequedJob.Image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

func updateClient(dequedJobUuid string, updateType int, message gin.H) {
	JobContextMap.Range(func(interfaceKey, interfaceVal interface{}) bool {
		// log.Println("\n\n\n\n\n\nranging over the contexts")
		uuid, keyOk := interfaceKey.(string)
		c, valOk := interfaceVal.(*gin.Context)
		if keyOk && valOk && subtle.ConstantTimeCompare([]byte(uuid), []byte(dequedJobUuid)) == 1 {
			var update gin.H
			switch updateType {
			case constants.Completed:
				update = gin.H{"message": "job completed", "code": 200}
				JobContextMap.Delete(uuid)

			case constants.Running:
				update = gin.H{"message": "job running", "code": 200}
			default:
				update = gin.H{"message": "job failed", "code": 401}
				JobContextMap.Delete(uuid)
			}
			sse.Encode(c.Writer, sse.Event{
				Event: "update",
				Data:  update,
			})

			return false
		}
		return true
	})
}
