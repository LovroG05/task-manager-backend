package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lovrog05/task-manager-backend/controllers"
	"github.com/lovrog05/task-manager-backend/models"
)

func main() {
	err := godotenv.Load("db.env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	r := gin.Default()

	models.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	r.GET("/users", controllers.FindUsers)
	r.POST("/users", controllers.CreateUser)
	r.GET("/users/:id", controllers.FindUser)
	r.PATCH("/users/:id", controllers.UpdateUser)
	r.DELETE("/users/:id", controllers.DeleteUser)

	r.Run()
}
