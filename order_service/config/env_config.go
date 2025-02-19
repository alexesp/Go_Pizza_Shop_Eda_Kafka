package config

import (
	"log"

	"github.com/joho/godotenv"
)

var env ConfigDto

var ConfigDto struct {
	port                string
	database_url        string
	database_name       string
	kafka_host          string
	kafka_port          string
	kafka_default_topic string
	kafka_group_id      string
}

func init() {
	if env.port == "" {
		ConfigEnv()
	}
}

func ConfigEnv() {
	LoadEnvVariable()
	env = ConfigDto{
	port                os.Getenv("PORT"),
	database_url        os.Getenv("MONGO_DB_URL"),
	database_name       os.Getenv("MONGO_DB_NAME"),
	kafka_host          os.Getenv("KAFKA_HOST"),
	kafka_port          os.Getenv("KAFKA_PORT"),
	kafka_default_topic os.Getenv("KAFKA_DEFAULT_TOPIC"),
	kafka_group_id      os.Getenv("KAFKA_GROUP_ID"),

	}
}



func LoadEnvVariable(){
	
	if _, err, os.Stat(".env"); err != nil{
		err := godotenv.Load()
		if err != nil{
			log.Fatalf("error loading .env file: %v", err)
		}
	}else {
		logge
	}
}