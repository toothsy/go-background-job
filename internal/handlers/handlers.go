package handlers

import (
	"fmt"
	"github/toothsy/go-background-job/helpers"
	"github/toothsy/go-background-job/internal/config"
	"github/toothsy/go-background-job/internal/constants"
	"github/toothsy/go-background-job/internal/models"
	"github/toothsy/go-background-job/internal/repository"
	dbrepo "github/toothsy/go-background-job/internal/repository/dbRepo"
	"github/toothsy/go-background-job/internal/workerpool"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
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
		DB:  dbrepo.NewMongoConnection(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Authenticate, allows the user to login
func Authenticate(c *gin.Context) {
	email := c.PostForm("email")
	passwordHash := c.PostForm("phash")
	//checking if the user exists in database and if they are verified or not
	u := &models.UserPayload{Email: email, PasswordHash: passwordHash, CreatedAt: time.Now()}
	databaseUser, err := Repo.DB.FetchUser(u)
	if err != nil {
		c.JSON(401, gin.H{"message": "user does not exist"})
		return
	}
	if !databaseUser.IsVerified {
		c.JSON(406, gin.H{"message": "user is not verified"})
		return
	}
	job := models.Job{
		Id:      uuid.New().String(),
		Status:  constants.Queued,
		JobType: constants.Authenticate,
		Image:   models.ImagePayload{},
		User:    *u,
	}
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Writer.Flush()

	workerpool.JobContextMap.Store(job.Id, c)
	Repo.App.WorkerPool.Enqueue(&job)
	sse.Encode(c.Writer, sse.Event{
		Event: "update",
		Data:  gin.H{"message": "job queued", "code": 200},
	})
	for {
		if _, ok := workerpool.JobContextMap.Load(job.Id); !ok {
			break
		}
	}
}

// UploadImage used to upload image
func UploadImage(c *gin.Context) {
	FILE_SIZE := 200 * 1024 //limit of 200KB
	username := c.Request.FormValue("username")
	u := &models.UserPayload{UserName: username}
	imageFile, err := c.FormFile("image")

	if err != nil {
		// Handle error
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	log.Println(imageFile.Filename, imageFile.Size)
	if imageFile.Size > int64(FILE_SIZE) {
		// Handle error
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "file size exceeds " + fmt.Sprintf("%d", FILE_SIZE/1000) + "KB limit"})
		return
	}
	i := &models.ImagePayload{Image: imageFile, UserName: *u, CreatedAt: time.Now()}
	job := models.Job{
		Id:      uuid.New().String(),
		Status:  constants.Queued,
		JobType: constants.Upload,
		Image:   *i,
		User:    *u,
	}
	Repo.App.WorkerPool.Enqueue(&job)
	c.JSON(http.StatusOK, gin.H{"message": "upload Job published to queue"})
}

func SignUp(c *gin.Context) {
	username := c.PostForm("username")
	passwordHash := c.PostForm("phash")
	email := c.PostForm("email")
	doubleHash, err := bcrypt.GenerateFromPassword([]byte(passwordHash), bcrypt.MinCost)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": 500, "error in hashing": err.Error()})
	}
	u := &models.UserPayload{UserName: username, PasswordHash: string(doubleHash), Email: email, CreatedAt: time.Now()}

	_, err = Repo.DB.FetchUser(u)
	// if user already exists then let frontend handle the redirect
	if err == nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": " user already exists in database", "code": 307})
		return
	}
	Repo.DB.SignUp(u)
	c.JSON(http.StatusOK, gin.H{"message": "inserted the data into database", "code": 200})
	verificationString := helpers.GenerateToken(u)
	mailData, err := helpers.ReadFile("./templates/minified.html")
	if err != nil {
		log.Println("err ", err)
	}
	redirectUrl := fmt.Sprintf("http://localhost:8080/projects/auth/verify?uuid=%s&uname=%s&eml=%s&", verificationString, u.UserName, u.Email)
	mailData = strings.Replace(mailData, "[%body%]", "Please click on the button below to verify your email", 1)
	mailData = strings.Replace(mailData, "[%url%]", redirectUrl, 1)
	mailData = strings.Replace(mailData, "[%button-content%]", "Click me to Verify", 1)
	mailData = strings.Replace(mailData, "[%sub-header%]", "aah General Kenobi", 1)
	from := mail.NewEmail("Go background Job Project", os.Getenv("FRMEML"))
	subject := "Hello please verify you email"
	to := mail.NewEmail(u.UserName, u.Email)
	message := mail.NewSingleEmail(from, subject, to, "", mailData)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRD"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(response.StatusCode)
		log.Println(response.Body)
		log.Println(response.Headers)
	}
}

func Verify(c *gin.Context) {
	verificationString := c.Query("uuid")
	uname := c.Query("uname")
	eml := c.Query("eml")
	user := &models.UserPayload{UserName: uname, Email: eml}
	if helpers.VerifyToken(verificationString, user) {
		Repo.DB.UpdateUserVerification(user)
		c.JSON(http.StatusAccepted, gin.H{"code": 200, "message": "all good you verified fam"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "oh NO somethings wrong"})

	}

}
