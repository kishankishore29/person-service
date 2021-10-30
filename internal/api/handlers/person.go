package handlers

import (
	"net/http"
	"person-service/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIServer struct {
	Database *gorm.DB
}

type ApiError struct {
	Error string `json:"error"`
}

// getPersons Returns all the person records stored in the database.
// TODO: Add pagination to this API
func (server APIServer) GetPersons(context *gin.Context) {
	var persons []models.Person

	// Query all the records present in the database
	result := server.Database.Find(&persons)

	// Check if there any error while query the records.
	if result.Error != nil {

		// Return a 500 Internal server error
		context.JSON(http.StatusInternalServerError, ApiError{Error: result.Error.Error()})
		return
	}

	// Return the serilized response with a 200 OK status code
	context.JSON(http.StatusOK, persons)
}

//createPerson Creates a new record in the person table
func (server APIServer) CreatePerson(context *gin.Context) {

	// Create a varialbe of the Person struct type
	var person models.Person

	// The request body will be unmarshalled to the pesron variable.
	err := context.BindJSON(&person)

	// Check if there was an error while unmarshalling the JSON request body.
	if err != nil {
		context.JSON(http.StatusBadRequest, ApiError{Error: err.Error()})
		return
	}

	//Create the new person record
	result := server.Database.Create(&person)

	// Check if there was an error while inserting the records to the database.
	if result.Error != nil {

		// Return a 500 Internal server error with the appropriate error message.
		context.JSON(http.StatusInternalServerError, ApiError{Error: result.Error.Error()})

	}

	// Return the person details.
	// The person variable now contains the latest details from the database
	context.JSON(http.StatusOK, person)

}
