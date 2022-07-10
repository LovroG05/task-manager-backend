package cron

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/go-co-op/gocron"
	"github.com/lovrog05/task-manager-backend/models"
	"google.golang.org/api/option"
)

var scheduler *gocron.Scheduler = gocron.NewScheduler(time.Local)

func StartScheduler() {
	scheduler.StartAsync()
}

func RegisterTaskCron(taskID uint) {
	var task models.Task
	if err := models.DB.Where("id = ?", taskID).Preload("Creator").Preload("Assignees").Find(&task).Error; err != nil {
		log.Println(err)
	}

	var days = []string{
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
		"Sunday",
	}

	var result []string
	for i := 0; i < len(days); i++ {
		if task.Days&(1<<uint(i)) != 0 {
			result = append(result, days[i])
		}
	}

	log.Println("Setting CRON for task " + task.Title + " on days " + strings.Join(result, ", ") + " at " + task.Time)

	for _, day := range result {
		switch strings.ToLower(day) {
		case "monday":
			scheduler.Every(1).Monday().At(task.Time).Do(makePushNotif, task)
		case "tuesday":
			scheduler.Every(1).Tuesday().At(task.Time).Do(makePushNotif, task)
		case "wednesday":
			scheduler.Every(1).Wednesday().At(task.Time).Do(makePushNotif, task)
		case "thursday":
			scheduler.Every(1).Thursday().At(task.Time).Do(makePushNotif, task)
		case "friday":
			scheduler.Every(1).Friday().At(task.Time).Do(makePushNotif, task)
		case "saturday":
			scheduler.Every(1).Saturday().At(task.Time).Do(makePushNotif, task)
		case "sunday":
			scheduler.Every(1).Sunday().At(task.Time).Do(makePushNotif, task)
		}
	}
}

func makePushNotif(task models.Task) {
	last := task.Last
	assignees := task.Assignees
	len_assignees := len(assignees)
	fmt.Println("last: ", last)
	fmt.Println("len: ", len_assignees)
	for assignee := range assignees {
		fmt.Println(assignees[assignee].Username)
	}

	if last == 0 {
		if len_assignees > 1 {
			task.Last = last + 1
			sendPushNotif(task, assignees[task.Last])
			models.DB.Save(&task)
		} else {
			sendPushNotif(task, assignees[task.Last])
		}
	} else {
		if len_assignees >= last+1 {
			task.Last = last + 1
			sendPushNotif(task, assignees[task.Last])
			models.DB.Save(&task)
		} else {
			task.Last = 0
			sendPushNotif(task, assignees[task.Last])
			models.DB.Save(&task)
		}
	}
}

func sendPushNotif(task models.Task, user models.User) {
	opt := option.WithCredentialsFile("servicekey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println(err)
		return
	}

	message := &messaging.Message{
		Token: user.FmcToken,
		Notification: &messaging.Notification{
			Title: "Task: " + task.Title,
			Body:  task.Description,
		},
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	response, err := client.Send(context.Background(), message)
	if err != nil {
		log.Println(err, response)
		return
	}
}
