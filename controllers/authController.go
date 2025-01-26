package controllers

import (
	"meeting-backend/models"
	"meeting-backend/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CurrentUser(c *gin.Context) {
	userID, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

func Login(c *gin.Context) {
	var input UserInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad login request"})
	}

	token, err := models.LoginCheck(input.Email, input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
