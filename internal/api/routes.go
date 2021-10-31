package api

import (
	"person-service/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func AddPersonRoutes(router *gin.RouterGroup, server handlers.APIServer) {
	personGroup := router.Group("/world/person/")
	{
		personGroup.GET("/", server.GetPersons)
		personGroup.POST("/", server.CreatePerson)
	}
}
