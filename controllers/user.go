package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/lovrog05/task-manager-backend/models"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User
	models.DB.Find(&users)
	c.JSON(200, gin.H{
		"message": "success",
		"data":    users,
	})
}
