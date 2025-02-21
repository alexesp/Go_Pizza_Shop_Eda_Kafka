package repository

import (
	"context"

	"github.com/alexesp/Go_Pizza_Shop_Eda_Kafka.git/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRepository interface {
	Create(data interface{}, ctx interface{}) (interface{}, error)
}

type MongoRepository struct {
	collection *mongo.Collection
}

func getSessionContext(sessionContext interface{}) mongo.SessionContext {
	cont := context.Background()
	if sessionContext == nil {
		return mongo.NewSessionContext(cont, mongo.SessionFromContext(cont))
	}

	return sessionContext.(mongo.SessionContext)
}

func (mr *MongoRepository) Create(data interface{}, ctx interface{}) (interface{}, error) {
	sc := getSessionContext(ctx)
	resut, err := mr.collection.InsertOne(sc, data)
	return resut, err
}

func GetMongoRepository(dbName, collectionName string) *MongoRepository {
	collection := config.GetDatabaseCollection(&dbName, collectionName)
	return &MongoRepository{
		collection: collection,
	}
}
