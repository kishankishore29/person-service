package internal

import (
	"fmt"
	"person-service/internal/models"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	Database *gorm.DB
	Router   *mux.Router
}

//InitializeDatabase Connect to the database and apply the migrations. Return the database connection handle.
func (server *Server) InitializeDatabase(databaseUser, databasePassword, databasePort, databaseHost, databaseName string) {

	logger, _ := zap.NewProduction()
	var err error
	
	// Create the database URL using the values passed to the startup function.
	databaseURL := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		databaseHost, databasePort, databaseUser, databaseName, databasePassword)

	// Open a connection to postgres
	server.Database, err = gorm.Open(postgres.Open(databaseURL))

	// Check if there was an error while opening a connection to the database
	if err != nil {
		logger.Error("There was an error while trying to connect to the database!")
		logger.Error(err.Error())
	}

	// Apply the migrations to the database.
	server.Database.Debug().AutoMigrate(&models.Person{})
}
