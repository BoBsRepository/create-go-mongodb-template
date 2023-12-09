package routes

import ("github.com/gin-gonic/gin"
		"gin-mongo-api/src/controllers"
		"gin-mongo-api/src/middlewear"
	)

func UserRoute(router *gin.Engine)  {

	router.GET("/",middleware.AuthMiddleware(),controllers.Greeting())
	apiGroup := router.Group("/api/auth")
    apiGroup.POST("/register", controllers.Register());
	apiGroup.POST("/login",controllers.Login())
};