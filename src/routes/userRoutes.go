package routes

import ("github.com/gin-gonic/gin"
		"gin-mongo-api/src/controllers")

func UserRoute(router *gin.Engine)  {
    router.POST("/user", controllers.Register())
};