package internal

import (
	"fmt"
	"person-service/internal/models"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//InitializeDatabase Connect to the database and apply the migrations. Return the database connection handle.
func InitializeDatabase(databaseUser, databasePassword, databasePort, databaseHost, databaseName string) *gorm.DB {

	logger, _ := zap.NewProduction()
	var err error

	// Create the database URL using the values passed to the startup function.
	databaseURL := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		databaseHost, databasePort, databaseUser, databaseName, databasePassword)

	// Open a connection to postgres
	database, err := gorm.Open(postgres.Open(databaseURL))

	// Check if there was an error while opening a connection to the database
	if err != nil {
		logger.Error("There was an error while trying to connect to the database!")
		logger.Error(err.Error())
	}

	// Apply the migrations to the database.
	database.Debug().AutoMigrate(&models.Person{})

	// Return the database object
	return database
}
