package repository

import "github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/config"

type Repositories struct {
	OrderRepository IRepository
}

func GetRepositories() *Repositories {
	return &Repositories{
		OrderRepository: GetMongoRepository(config.GetEnvProperty("database_name"), "orders"),
	}
}
