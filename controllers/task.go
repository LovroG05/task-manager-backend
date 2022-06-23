package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lovrog05/task-manager-backend/models"
	"github.com/lovrog05/task-manager-backend/utils/cron"
	"github.com/lovrog05/task-manager-backend/utils/token"
)

func GetTasks(c *gin.Context) {
	var tasks []models.Task
	models.DB.Preload("Creator").Preload("Assignees").Find(&tasks)
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": tasks})
}

func FindTask(c *gin.Context) {
	var tasks []models.Task
	if err := models.DB.Where("id = ?", c.Param("id")).Preload("Creator").Preload("Assignees").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": tasks})
}

type TaskInput struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description" binding:"required"`
	Daily        bool   `json:"daily" binding:"required"`
	OnTime       string `json:"on_time" binding:"required"`
	OnDay        string `json:"on_day"`
	AssigneesIDs []uint `json:"assignees_ids"`
}

func CreateTask(c *gin.Context) {
	var task TaskInput
	if err := c.ShouldBindJSON(&task); err != nil {
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
	//TODO: check if time is valid
	fmt.Println("assignees: ", task.AssigneesIDs)
	var assignees []models.User
	for id := range task.AssigneesIDs {
		assignee, err := models.GetUserByID(task.AssigneesIDs[id])
		if err == nil {
			assignees = append(assignees, assignee)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	fmt.Println("assignees: ", assignees)

	newtask := models.Task{
		Title:       task.Title,
		Description: task.Description,
		Daily:       task.Daily,
		OnTime:      task.OnTime,
		OnDay:       task.OnDay,
		Assignees:   assignees,
		Creator:     u,
		Last:        0,
	}

	err = models.DB.Create(&newtask).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cron.RegisterTaskCron(newtask.ID)
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": newtask})
}
