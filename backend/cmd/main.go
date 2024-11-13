package main

import (
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/lunatictiol/referal-system/config"
	"github.com/lunatictiol/referal-system/db"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {

	envFile, err := godotenv.Read(".env")
	if err != nil {
		config.LogFatal("cannot load env")
	}

	var apiServer ApiServer
	dbburl := envFile["MONGO_URL"]
	db, err := db.NewMongoDbConnection(dbburl)
	if err != nil {
		config.LogFatal("error connecting to db")
		return
	}
	ping(db)
	apiServer.New("8080", db)
	config.Logger.Info("server is starting at", "port", "8080")
	err = apiServer.Run()

	if err != nil {
		config.LogFatal("error running api", err)
	}

}

func ping(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	err := client.Ping(ctx, readpref.Primary())
	if err != nil {

		config.LogFatal("cannot ping mongodb", err)
	}
}
