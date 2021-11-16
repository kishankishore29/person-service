package api

import (
	"person-service/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

// AddPersonRoutes Register the routes with the corresponding functions. The URLs are grouped by their path.
func AddPersonRoutes(router *gin.RouterGroup, server handlers.APIServer) {
	personGroup := router.Group("/world/person/")
	{
		personGroup.GET("/", server.GetPersons)
		personGroup.POST("/", server.CreatePerson)

		personGroup.GET("/:personId", server.GetPerson)
		personGroup.PUT("/:personId", server.UpdatePerson)
		personGroup.DELETE("/:personId", server.DeletePerson)
		personGroup.PATCH("/:personId", server.PartialUpdatePerson)

	}
}
