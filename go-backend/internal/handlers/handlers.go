package handlers

import (
	"fmt"
	"github/toothsy/go-background-job/internal/config"
	"github/toothsy/go-background-job/internal/constants"
	"github/toothsy/go-background-job/internal/models"
	"github/toothsy/go-background-job/internal/repository"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {

	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Authenticate, allows the user to login
func Authenticate(c *gin.Context) {
	username := c.PostForm("username")
	passwordHash := c.PostForm("phash")
	log.Println("username \t", username, "and password", passwordHash)
	job := models.Job{
		Id:           uuid.New().String(),
		Status:       constants.Queued,
		JobType:      constants.Authenticate,
		Image:        nil,
		Username:     username,
		PasswordHash: passwordHash,
	}
	Repo.App.WorkerPool.Enqueue(job)
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"message": "auth Job published to queue"})
}

// UploadImage used to upload image
func UploadImage(c *gin.Context) {
	FILE_SIZE := 200 * 1024 //limit of 200KB
	c.Header("Content-Type", "application/json")
	username := c.Request.FormValue("username")
	imageFile, err := c.FormFile("image")
	if err != nil {
		// Handle error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(imageFile.Filename, imageFile.Size)
	if imageFile.Size > int64(FILE_SIZE) {
		// Handle error
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds " + fmt.Sprintf("%d", FILE_SIZE/1000) + "KB limit"})
		return
	}
	job := models.Job{
		Id:       uuid.New().String(),
		Status:   constants.Queued,
		JobType:  constants.Upload,
		Image:    imageFile,
		Username: username,
	}
	Repo.App.WorkerPool.Enqueue(job)
	c.JSON(http.StatusOK, gin.H{"message": "upload Job published to queue"})
}

func SignUp(c *gin.Context) {
	// username := c.PostForm("username")
	// passwordHash := c.PostForm("phash")

}
