package router

import (
	"github.com/gin-gonic/gin"
	"main.go/controllers"
)

func authRouter() {
	router := gin.Default()
	router.GET("/auth/getUsers", controllers.getUser)
	router.POST("/auth/login", controllers.loginUser)
	router.Run(":8080")
}
