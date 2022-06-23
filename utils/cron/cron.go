package cron

import (
	"context"
	"fmt"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/go-co-op/gocron"
	"github.com/lovrog05/task-manager-backend/models"
	"google.golang.org/api/option"
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
	opt := option.WithCredentialsFile("servicekey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Println(err)
		return
	}

	message := &messaging.Message{
		Token: user.FmcToken,
		Notification: &messaging.Notification{
			Title: "Task: " + title,
			Body:  "do your task now",
		},
		Data: map[string]string{"title": title, "body": "do your task now"},
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	time.Sleep(10 * time.Second) // sleep for 10 seconds

	response, err := client.Send(context.Background(), message)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(response)
		return
	}
}
