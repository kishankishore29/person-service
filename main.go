package main

import (
	"person-service/config"
	"person-service/internal"
	"person-service/internal/api"
	"person-service/internal/api/handlers"
	"person-service/internal/seed"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()

	// Load the configuration to be used by the application
	configuration, err := config.LoadConfig("local")

	if err != nil {
		logger.Error("There was a problem while reading the configuration!")
		logger.Error(err.Error())
		return
	}

	// Initialize the database. Pass all the values required. This will also apply the migrations.
	database := internal.InitializeDatabase(configuration.DatabaseUser, configuration.DatabasePassword,
		configuration.DatabasePort, configuration.DatabaseHost, configuration.DatabaseName)

	// Check if data needs to be seeded
	if configuration.ShouldSeedData {
		// Initialize the database with random entries of persons.
		seed.LoadRandomPersonData(int32(configuration.NumberOfTestPersonEntries), database)
	}

	// Create an instance of APIServer struct and set the database field
	apiServer := handlers.APIServer{Database: database}

	// Create a HTTP router
	router := gin.Default()

	// Create a new v1 group.
	// This will help differentiate between the further versions of ther API.
	v1Router := router.Group("/v1")

	// Call the function to add the routes related to the Person resource.
	api.AddPersonRoutes(v1Router, apiServer)

	// Start the HTTP server by passing the URL that it needs to start on.
	router.Run(":" + configuration.HTTPPort)
}
