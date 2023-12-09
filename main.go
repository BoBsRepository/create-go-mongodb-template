package main

import (
	"gin-mongo-api/src/database"
	"gin-mongo-api/src/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	database.ConnectDB()
	routes.UserRoute(router)
	router.Run(":6000")
}
