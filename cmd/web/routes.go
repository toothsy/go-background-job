package main

import (
	"encoding/json"
	"github/toothsy/go-background-job/internal/config"
	"github/toothsy/go-background-job/internal/handlers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func routes(app *config.AppConfig) http.Handler {
	mux := gin.Default()
	mux.Use(cors.Default())
	mux.Use(gin.Recovery())
	// mux.Use(JSONmiddleware())
	mux.Use(increaseBufferMiddleware)
	auth := mux.Group("/projects/auth")
	{
		auth.POST("/login", handlers.Authenticate)
		auth.POST("/signup", handlers.SignUp)
		auth.GET("/verify", handlers.Verify)
	}
	mux.POST("/projects/upload/", handlers.UploadImage)

	return mux
}

// returns stuf in JSON format
func JSONmiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			log["remote_addr"] = params.ClientIP
			log["response_time"] = params.Latency.String()

			s, _ := json.MarshalIndent(log, "", "\n")
			return string(s) + "\n"
		},
	)
}

// increaseBufferMiddleware, increases the buffer size to 1MB so I can possible get 1mb imaages and reject
func increaseBufferMiddleware(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
	c.Next()
}
