package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lovrog05/task-manager-backend/controllers"
	"github.com/lovrog05/task-manager-backend/middlewares"
	"github.com/lovrog05/task-manager-backend/models"
)

func main() {
	err := godotenv.Load("db.env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	models.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	r := gin.Default()
	public := r.Group("/api")

	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)

	r.Run()
}
