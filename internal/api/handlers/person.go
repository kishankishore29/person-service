package handlers

import (
	"errors"
	"net/http"
	"person-service/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	// Return the serialized response with a 200 OK status code
	context.JSON(http.StatusOK, persons)

}

//createPerson Creates a new record in the person table
func (server APIServer) CreatePerson(context *gin.Context) {

	// Create a varialbe of the Person struct type
	var person models.Person

	// The request body will be unmarshalled to the person variable.
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

		return

	}

	// Return the person details. The status code is 201 Created and not 200 OK.
	// The person variable now contains the latest details from the database
	context.JSON(http.StatusCreated, person)

}

// GetPerson Returns a person resource corresponding to the passed id.
func (server APIServer) GetPerson(context *gin.Context) {

	// This will store the result.
	var person models.Person

	// Check the id string passed is a valid uuid v4.
	if !validateUUID(context.Param("personId")) {
		context.JSON(http.StatusBadRequest, ApiError{Error: "Invalid uuid!"})
		return
	}

	// Get the person using the personId query parameter
	result := server.Database.First(&person, "id = ?", context.Param("personId"))

	if result.Error != nil {

		// Return a 404 not found if the id was found in the database.
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, ApiError{Error: result.Error.Error()})
			return
		}

		// Return a 500 Internal server error
		context.JSON(http.StatusInternalServerError, ApiError{Error: result.Error.Error()})
		return
	}

	// Return the serialized response with a 200 OK status code
	context.JSON(http.StatusOK, person)
}

// UpdatePerson Updates the person resource with the given id or else returns a 404 status code.
func (server APIServer) UpdatePerson(context *gin.Context) {

	// Check the id string passed is a valid uuid v4.
	if !validateUUID(context.Param("personId")) {
		context.JSON(http.StatusBadRequest, ApiError{Error: "Invalid uuid!"})
		return
	}

	// This will store the request body.
	var person models.Person

	// Get the person using the personId query parameter. This to check if the id exists in the database
	result := server.Database.First(&person, "id = ?", context.Param("personId"))

	// Return a 404 not found if the id was found in the database.
	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, ApiError{Error: result.Error.Error()})
			return
		}
	}

	// reset the values received in the fetch call.
	person = models.Person{}

	context.BindJSON(&person)

	person.Id = context.Param("personId")

	// Update the columns of the record.
	result = server.Database.Save(&person)

	if result.Error != nil {

		// Return a 404 not found if the id was found in the database.
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, ApiError{Error: result.Error.Error()})
			return
		}

		// Return a 500 Internal server error
		context.JSON(http.StatusInternalServerError, ApiError{Error: result.Error.Error()})
		return
	}

	// Return the serialized response with a 200 OK status code
	context.JSON(http.StatusOK, person)
}

// PartialUpdatePerson Partially updates the person resource with the given id or else returns a 404 status code.
func (server APIServer) PartialUpdatePerson(context *gin.Context) {

	// Check the id string passed is a valid uuid v4.
	if !validateUUID(context.Param("personId")) {
		context.JSON(http.StatusBadRequest, ApiError{Error: "Invalid uuid!"})
		return
	}

	// This will store the request body.
	var person models.Person

	context.BindJSON(&person)

	person.Id = context.Param("personId")

	// Update the columns of the record and retrun all the columns after the update is finished.
	result := server.Database.Model(&person).Clauses(clause.Returning{}).Updates(&person)

	// Return a 404 not found if the id was found in the database.
	if result.RowsAffected == 0 || result.Error != nil {

		if result.Error != nil {
			// Return a 500 Internal server error
			context.JSON(http.StatusInternalServerError, ApiError{Error: result.Error.Error()})
			return
		}

		context.JSON(http.StatusNotFound, ApiError{Error: "Record Not Found!"})
		return
	}

	// Get the person using the personId query parameter to get the latest values.
	server.Database.First(&person, "id = ?", context.Param("personId"))

	// Return the serialized response with a 200 OK status code
	context.JSON(http.StatusOK, person)
}

// DeletePerson Deletes the person record with the given id.
func (server APIServer) DeletePerson(context *gin.Context) {

	// Check the id string passed is a valid uuid v4.
	if !validateUUID(context.Param("personId")) {
		context.JSON(http.StatusBadRequest, ApiError{Error: "Invalid uuid!"})
		return
	}

	// Delete the record using the passed id.
	result := server.Database.Delete(&models.Person{}, "id = ?", context.Param("personId"))

	if result.Error != nil {

		// Return a 404 not found if the id was found in the database.
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, ApiError{Error: result.Error.Error()})
			return
		}

		// Return a 500 Internal server error
		context.JSON(http.StatusInternalServerError, ApiError{Error: result.Error.Error()})
		return
	}

	// Return the serialized response with a 200 OK status code
	context.JSON(http.StatusNoContent, nil)
}
