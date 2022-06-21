package cron

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/lovrog05/task-manager-backend/models"
)

func RegisterTaskCron(taskID uint) {
	var task models.Task
	if err := models.DB.Where("id = ?", taskID).Preload("Creator").Preload("Assignees").Find(&task).Error; err != nil {
		fmt.Println(err)
	}
	scheduler := gocron.NewScheduler(time.Local)

	if task.Daily {
		scheduler.Every(1).Day().At(task.OnTime).Do(func() {
			makePushNotif(taskID)
		})
	} else {
		if task.OnDay == "MONDAY" {
			scheduler.Every(1).Monday().At(task.OnTime).Do(func() {
				makePushNotif(taskID)
			})
		} else if task.OnDay == "TUESDAY" {
			scheduler.Every(1).Tuesday().At(task.OnTime).Do(func() {
				makePushNotif(taskID)
			})
		} else if task.OnDay == "WEDNESDAY" {
			scheduler.Every(1).Wednesday().At(task.OnTime).Do(func() {
				makePushNotif(taskID)
			})
		} else if task.OnDay == "THURSDAY" {
			scheduler.Every(1).Thursday().At(task.OnTime).Do(func() {
				makePushNotif(taskID)
			})
		} else if task.OnDay == "FRIDAY" {
			scheduler.Every(1).Friday().At(task.OnTime).Do(func() {
				makePushNotif(taskID)
			})
		} else if task.OnDay == "SATURDAY" {
			scheduler.Every(1).Saturday().At(task.OnTime).Do(func() {
				makePushNotif(taskID)
			})
		} else if task.OnDay == "SUNDAY" {
			scheduler.Every(1).Sunday().At(task.OnTime).Do(func() {
				makePushNotif(taskID)
			})
		}
	}
}

func makePushNotif(taskID uint) {
	var task models.Task
	if err := models.DB.Where("id = ?", taskID).Preload("Creator").Preload("Assignees").Find(&task).Error; err != nil {
		fmt.Println(err)
	}
	last := task.Last
	assignees := task.Assignees
	len_assignees := len(assignees)
	if len_assignees >= last+1 {
		task.Last = last + 1
		sendPushNotif(task.Title, assignees[task.Last])
		models.DB.Save(&task)
	} else {
		task.Last = 0
		sendPushNotif(task.Title, assignees[task.Last])
		models.DB.Save(&task)
	}
}

func sendPushNotif(title string, user models.User) {
	fmt.Println("Task " + title + " must be performed by " + user.Username + " right now")
}
