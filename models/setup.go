package models

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(host string, port string, user string, password string, dbname string) {
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable TimeZone=Europe/Berlin"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Panic("Failed to connect to database: ", err)
	}

	database.AutoMigrate(&User{})
	database.AutoMigrate(&Task{})
	database.AutoMigrate(&TaskLog{})

	DB = database
}
