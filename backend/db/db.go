package db

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoDbConnection(connStr string) (*mongo.Client, error) {

	client, err := mongo.Connect(options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, err
	}
	return client, nil

}
