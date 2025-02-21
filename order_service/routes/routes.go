package routes

import (
	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, publisher service.IMessagePublisher) {
	orderRoutes := r.Group("/order-service")
	{
		registerOrderRoutes(orderRoutes, publisher)
	}
}
