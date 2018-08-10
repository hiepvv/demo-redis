package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hiepvv/auth-no-redis/models"
	"github.com/hiepvv/auth-no-redis/routers"
)

var router *gin.Engine

func init() {
	models.InitialDBSession()

	router = gin.Default()
}

func main() {
	api := router.Group("/api")
	{
		api.POST("users/register", routers.AddMulUserEndPoint)
		api.GET("/users", routers.FindMul)
		api.POST("/remove", routers.EraseRedis)
	}

	router.Run(":3000")
}
