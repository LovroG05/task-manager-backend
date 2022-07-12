package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lovrog05/task-manager-backend/models"
	"github.com/lovrog05/task-manager-backend/utils"
)

func GetAllUsers(c *gin.Context) {
	pagination := utils.NewPagination(c)

	var users []models.User
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := models.DB.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data":    users,
	})
}
