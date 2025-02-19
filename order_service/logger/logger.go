package logger

import (
	"os"
)

func Log(message any) {
	isLoggedEnabled := os.Getenv("log")

	if isLoggedEnabled != "" {
		utils.AppendToFile("log.txt", message.(string))
	}
}
