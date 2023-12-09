package routes

import ("github.com/gin-gonic/gin"
		"gin-mongo-api/src/controllers")

func UserRoute(router *gin.Engine)  {
	router.GET("/",controllers.Greeting())
    router.POST("/register", controllers.Register());
	router.POST("/login",controllers.Login())
};