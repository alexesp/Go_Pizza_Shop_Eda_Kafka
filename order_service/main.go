package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()

	app.Use(gin.Recovery())

	app.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "service is up and runnig",
		})
	})

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", 8001),
		Handler:        app,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB

	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
