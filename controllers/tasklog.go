package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lovrog05/task-manager-backend/models"
	"github.com/lovrog05/task-manager-backend/models/inputs"
	"github.com/lovrog05/task-manager-backend/utils"
	"github.com/lovrog05/task-manager-backend/utils/token"
)

func LogTask(c *gin.Context) {
	var tasklog inputs.TaskLogInput
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

func GetTaskLogsByTask(c *gin.Context) {
	task_id := c.Param("task_id")

	pagination := utils.NewPagination(c)
	// str to int
	task_id_int, err := strconv.Atoi(task_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err := findTask(task_id_int)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var logs []models.TaskLog

	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := models.DB.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.Where("task_id = ?", task.ID).Find(&logs)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": logs})
}

func GetTaskLogsByUser(c *gin.Context) {
	pagination := utils.NewPagination(c)

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

	var logs []models.TaskLog
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := models.DB.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.Where("user = ?", u).Find(&logs)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": logs})
}
