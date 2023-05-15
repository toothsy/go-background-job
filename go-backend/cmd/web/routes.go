package main

import (
	"github/toothsy/go-background-job/internal/config"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func routes(app *config.AppConfig) http.Handler {
	mux := gin.Default()
	mux.Use(cors.Default())
	mux.GET("/", func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		ctx.String(http.StatusOK, "Hello World from root path")
	})
	auth := mux.Group("/auth/")
	{
		auth.GET("/", func(ctx *gin.Context) {
			ctx.Header("Content-Type", "application/json")

			ctx.String(http.StatusOK, "Hello World from auth root")
		})
		auth.GET("/login", func(ctx *gin.Context) {
			ctx.Header("Content-Type", "application/json")

			ctx.String(http.StatusOK, "Hello World from login")
		})

	}

	return mux
}
