package internal

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Provides the database and the router fields.
type Server struct {
	Database *gorm.DB
	Router   *gin.Engine
}

func (server *Server) Run(address string) {

	// Run the gin http server
	server.Router.Run(address)
}
