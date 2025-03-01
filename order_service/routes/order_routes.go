package routes

import (
	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/handler"
	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/service"
	"github.com/gin-gonic/gin"
)

func registerOrderRoutes(r *gin.RouterGroup, publisher service.IMessagePublisher) {
	orderHandler := handler.GetOrderHandler(publisher)

	r.POST(
		"/create",
		orderHandler.CreateOrder,
	)
}
