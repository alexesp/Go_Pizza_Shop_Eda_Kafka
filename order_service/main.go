package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/constants"
	messageconsumer "github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/message-consumer"
	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/repository"
	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/routes"
	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/service"
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

	var repositories = repository.GetRepositories()
	var orderConsumer = messageconsumer.GetOrderMessageConsumer(
		service.GetNewKafkaConsumer(constants.TOPIC_ORDER, "order-message"),
		*repositories,
	)

	routes.RegisterRoutes(
		app,
		service.GetKafkaMessagePublisher(constants.TOPIC_ORDER),
	)

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
