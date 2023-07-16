package util

import (
	"di/model"
	"fmt"
	"log"
	"os"
	"reflect"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() (*gorm.DB, error) {

	pgHost, exists := os.LookupEnv("POSTGRES_HOST")

	if !exists {
		panic("POSTGRES_HOST is not defined!")
	}
	pgPort, exists := os.LookupEnv("POSTGRES_PORT")

	if !exists {
		panic("POSTGRES_PORT is not defined!")
	}

	pgUser, exists := os.LookupEnv("POSTGRES_USER")

	if !exists {
		panic("POSTGRES_USER is not defined!")
	}

	pgPassword, exists := os.LookupEnv("POSTGRES_PASSWORD")

	if !exists {
		panic("POSTGRES_PASSWORD is not defined!")
	}

	pgDB, exists := os.LookupEnv("POSTGRES_DB")

	if !exists {
		panic("POSTGRES_DB is not defined!")
	}

	pgSSL, exists := os.LookupEnv("POSTGRES_SSL")

	if !exists {
		panic("POSTGRES_SSL is not defined!")
	}

	pgTimezone, exists := os.LookupEnv("POSTGRES_TIMEZONE")

	if !exists {
		panic("POSTGRES_TIMEZONE is not defined!")
	}

	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s", pgHost, pgPort, pgUser, pgPassword, pgDB, pgSSL, pgTimezone)

	gormDB, err := gorm.Open(postgres.Open(pgConnString), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

func CreateOrUpdateSchema(db *gorm.DB) error {

	// add models here
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalln(err)
		return err
	}

	if err := db.AutoMigrate(&model.Pipeline{}); err != nil {
		log.Fatalln(err)
		return err
	}

	if err := db.AutoMigrate(&model.PipelineSchedule{}); err != nil {
		log.Fatalln(err)
		return err
	}

	if err := db.AutoMigrate(&model.RunStatus{}); err != nil {
		log.Fatalln(err)
		return err
	}

	if err := createDefaultRunStatuses(db); err != nil {
		log.Fatalln(err)
		return err
	}

	if err := db.AutoMigrate(&model.Run{}); err != nil {
		log.Fatalln(err)
		return err
	}

	if err := db.AutoMigrate(&model.RunStepStatus{}); err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

func createDefaultRunStatuses(db *gorm.DB) error {
	defaultRunStatuses := []*model.RunStatus{
		{Name: "Not Run", IsFinal: false},
		{Name: "Executing", IsFinal: false},
		{Name: "Error", IsFinal: true},
		{Name: "Success", IsFinal: true},
	}

	for index, status := range defaultRunStatuses {
		str := reflect.ValueOf(status).Elem()

		if str.Kind() == reflect.Struct {
			pipelineIDField := str.FieldByName("ID")
			if pipelineIDField.IsValid() {
				if pipelineIDField.CanSet() {
					if pipelineIDField.Kind() == reflect.Uint {
						pipelineIDField.SetUint(uint64(index + 1))
					}
				}
			}
		}

		if result := db.FirstOrCreate(status, status); result.Error != nil {
			log.Fatalln(result.Error)
			return result.Error
		}
	}

	return nil
}
