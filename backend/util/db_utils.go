package util

import (
	"di/model"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() *gorm.DB {
	// load env variables
	pgHost := os.Getenv("POSTGRES_HOST")
	// pgPort := os.Getenv("POSTGRES_PORT")
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDB := os.Getenv("POSTGRES_DB")
	pgSSL := os.Getenv("POSTGRES_SSL")
	pgTimezone := os.Getenv("POSTGRES_TIMEZONE")

	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s", pgHost, "5432", pgUser, pgPassword, pgDB, pgSSL, pgTimezone)

	gormDB, err := gorm.Open(postgres.Open(pgConnString), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	return gormDB
}

func CreateOrUpdateSchema(db *gorm.DB) error {

	// add models here
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalln(err)
	}

	return nil
}
