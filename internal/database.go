package internal

import (
	"fmt"
	"person-service/internal/models"

	"github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

//InitializeDatabase Connect to the database and apply the migrations. Return the database connection handle.
func InitializeDatabase(url string) (*gorm.DB, error) {

	logger, _ := zap.NewProduction()
	var err error

	// Extract the database configuration using the DATABASE_URL environment variable.
	parsedURL, err := pq.ParseURL(url)

	if err != nil {
		logger.Error(fmt.Sprintf("There was an error while parsing the postgres URL : %e", err))
		return nil, err
	}

	// Open a connection to postgres
	database, err := gorm.Open(postgres.Open(parsedURL), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})

	// Check if there was an error while opening a connection to the database
	if err != nil {
		logger.Error(fmt.Sprintf("There was an error while connecting to the database : %e", err))
		return nil, err
	}

	// Apply the migrations to the database.
	err = database.Debug().AutoMigrate(&models.Person{})

	if err != nil {
		logger.Error(fmt.Sprintf("There was an error while applying migrations : %e", err))
		return nil, err
	}

	// Return the database object
	return database, nil
}
