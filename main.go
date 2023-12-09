package main

import (
	"gin-mongo-api/src/database"
	"gin-mongo-api/src/routes"
	"log"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	router.Use(limits.RequestSizeLimiter(100))
	
	database.ConnectDB(); 
	router.Use(secure.New(secure.Config{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	}))
	routes.UserRoute(router)

	if err := router.Run(":6000"); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
