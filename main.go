package main

import (
	"person-service/config"
	"person-service/internal"
	"person-service/internal/seed"

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

	// Create an instance of the Server struct and initialize the database
	server := internal.Server{}

	// Initialize the database. Pass all the values required. This will also apply the migrations.
	server.InitializeDatabase(configuration.DatabaseUser, configuration.DatabasePassword, configuration.DatabasePort, configuration.DatabaseHost, configuration.DatabaseName)

	if configuration.ShouldSeedData {
		// Initialize the database with random entries of persons.
		seed.LoadRandomPersonData(int32(configuration.NumberOfTestPersonEntries), server.Database)
	}

}
