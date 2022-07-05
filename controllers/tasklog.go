package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lovrog05/task-manager-backend/models"
	"github.com/lovrog05/task-manager-backend/utils/token"
)

type TaskLogInput struct {
	TaskID int    `json:"task_id" binding:"required"`
	Time   string `json:"time"`
	Notes  string `json:"notes"`
}

func LogTask(c *gin.Context) {
	var tasklog TaskLogInput
	if err := c.ShouldBindJSON(&tasklog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := findTask(tasklog.TaskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newlog := models.TaskLog{
		Task:  *task,
		Time:  tasklog.Time,
		Notes: tasklog.Notes,
		User:  u,
	}

	err = models.DB.Create(&newlog).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println(err.Error())
		return
	}

	//TODO change last in task record
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": newlog})
}
