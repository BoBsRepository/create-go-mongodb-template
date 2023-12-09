package main

import (
	"gin-mongo-api/src/database"
	"gin-mongo-api/src/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://google.com"}
	router.Use(limits.RequestSizeLimiter(10))
	database.ConnectDB()
	router.Use(secure.New(secure.Config{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	}))
	routes.UserRoute(router)
	router.Run(":6000")
}
