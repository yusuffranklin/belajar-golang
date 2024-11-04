package routes

import (
	"net/http"
	"rest-api-practice/models"
	"rest-api-practice/utils"

	"github.com/gin-gonic/gin"
)

func createUser(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "couldn't parse the requested data."})
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "couldn't create user."})
	}

	context.JSON(http.StatusOK, gin.H{"message": "user created successfully."})
}

func login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "couldn't parse the requested data"})
		return
	}

	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "couldn't authenticate user"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "login succeeded", "token": token})
}
