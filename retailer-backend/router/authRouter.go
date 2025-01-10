package router

import (
	"github.com/gin-gonic/gin"
)

func authRouter() {
	router := gin.Default()
	router.GET("/auth/getUsers", getUsers)
}
