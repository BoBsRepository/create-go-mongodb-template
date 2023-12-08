package main

import (
	"gin-mongo-api/src/database"
	"github.com/gin-gonic/gin"
	"gin-mongo-api/src/routes"
)

func main() {
	router := gin.Default()
	database.ConnectDB()
	routes.UserRoute(router)
	router.Run(":6000")
}
