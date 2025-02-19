package logger

import (
	"os"

	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka/order_service/utils"
)

func Log(message any) {
	isLoggedEnabled := os.Getenv("log")

	if isLoggedEnabled != "" {
		utils.AppendToFile("log.txt", message.(string))
	}
}
