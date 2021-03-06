package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lovrog05/task-manager-backend/controllers"
	"github.com/lovrog05/task-manager-backend/middlewares"
	"github.com/lovrog05/task-manager-backend/models"
	"github.com/lovrog05/task-manager-backend/utils/cron"
)

func main() {
	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	log.SetOutput(file)
	log.Println("Starting server...")

	err = godotenv.Load("db.env")
	if err != nil {
		log.Fatalln("Error loading .env file: ", err)
	}

	models.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	r := gin.Default()
	public := r.Group("/api")

	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)
	protected.GET("/allusers", controllers.GetAllUsers)
	protected.GET("/tasks", controllers.GetTasks)
	protected.POST("/tasks", controllers.CreateTask)
	protected.GET("/task/:id", controllers.FindTask)
	protected.DELETE("/task/delete/:id", controllers.DeleteTask)
	protected.PATCH("/user/updatefmc", controllers.UpdateFmcToken)
	protected.POST("/task/log", controllers.LogTask)
	protected.GET("/task/logs/:id", controllers.GetTaskLogsByTask)
	protected.GET("/task/logs", controllers.GetTaskLogsByUser)

	cron.StartScheduler()
	r.Run()
}
