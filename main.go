package main

import (
	"person-service/config"
	"person-service/internal"
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

	// Create an instance of the Server struct and initialize the database.
	// We will use the gin default router
	server := internal.Server{Router: gin.Default()}

	// Initialize the database. Pass all the values required. This will also apply the migrations.
	server.InitializeDatabase(configuration.DatabaseUser, configuration.DatabasePassword, configuration.DatabasePort, configuration.DatabaseHost, configuration.DatabaseName)

	// Check if data needs to be seeded
	if configuration.ShouldSeedData {
		// Initialize the database with random entries of persons.
		seed.LoadRandomPersonData(int32(configuration.NumberOfTestPersonEntries), server.Database)
	}

	// Create an instance of APIServer struct and set the database field
	apiServer := handlers.APIServer{Database: server.Database}

	// Register the GET route to the default router.
	// TODO Create a different fiole for storing the routes called routes.go.
	server.Router.GET("/v1/world/person/", apiServer.GetPersons)
	server.Router.POST("/v1/world/person/", apiServer.CreatePerson)

	// Run the HTTP server pasaing the address to it
	server.Run(configuration.HTTPHost + ":" + configuration.HTTPPort)
}
