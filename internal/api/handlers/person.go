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
	}

	// Return the serilized response with a 200 OK status code
	context.JSON(http.StatusOK, persons)
}
