package main

import (
	"github/toothsy/go-background-job/internal/config"
	"github/toothsy/go-background-job/internal/handlers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func routes(app *config.AppConfig) http.Handler {

	mux := gin.Default()
	mux.Use(cors.Default())
	mux.Use(gin.Logger())
	mux.Use(increaseBufferMiddleware)
	auth := mux.Group("/auth")
	{
		auth.POST("/login", handlers.Authenticate)

	}
	mux.POST("/upload/", handlers.UploadImage)

	return mux
}

// increaseBufferMiddleware, increases the buffer size to 1MB so I can possible get 1mb imaages and reject
func increaseBufferMiddleware(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
	c.Next()
}
