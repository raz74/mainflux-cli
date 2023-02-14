package repository

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"hamta-sharif/models"
	"log"
	"os"
)

func Initialize() *gorm.DB {
	dsn := getConfig()
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect database!")
	}
	err = database.AutoMigrate(&models.File{})
	if err != nil {
		log.Fatal(err)
	}

	return database
}

func getConfig() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	host := os.Getenv("DATABASE_HOST")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	port := os.Getenv("PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, dbname, port)
	return dsn
}
