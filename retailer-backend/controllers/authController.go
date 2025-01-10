package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/models"
)

func getUsers(context *gin.Context, users *models.User) {
	context.IndentedJSON(http.StatusOK, users)
}
func loginUser(context *gin.Context, users *models.User) {

}

func registerUsers() {

}
func forgotPassword() {

}
