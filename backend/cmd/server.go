package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lunatictiol/referal-system/services/users"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ApiServer struct {
	address string
	store   *mongo.Client
}

func (a *ApiServer) New(address string, store *mongo.Client) {
	a.address = address
	a.store = store

}

func (a *ApiServer) Run() error {
	router := gin.Default()

	gin.SetMode(gin.TestMode)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "working",
		})
	})
	api := router.Group("api/v1")
	userStore := users.NewStore(a.store)

	userRoutes := users.NewHandler(*userStore)
	userRoutes.RegisterRoutes(api)

	return router.Run(fmt.Sprintf(":%s", a.address))
}
