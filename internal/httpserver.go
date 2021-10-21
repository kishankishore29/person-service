package internal

import (
	"fmt"
	"log"
	"net/http"
)

func (server *Server) Run(address string) {
	fmt.Printf("Starting HTTP server on address: %s", address)

	// Start the HTTP server on the passed address.
	err := http.ListenAndServe(address, server.Router)

	// Check if there was an error while starting the HTTP server.
	if err != nil {
		log.Fatal("There was a problem while starting the HTTP server!")
	}
}
