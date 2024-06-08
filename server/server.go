package server

import (
	"RyanFin/GoPaseto/token"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	tokenMaker *token.PasetoMaker
	router     *gin.Engine
}

// @ Function creates a new server instance and returns it with the configured token maker generator
func NewServer(address string) (*Server, error) {
	// Initialize the PASETO token maker
	tokenMaker, err := token.NewPaseto("abcdefghijkl12345678901234567890")
	if err != nil {
		return nil, fmt.Errorf("Could not create token maker: %w", err)
	}

	// Create a new server instance
	server := &Server{
		tokenMaker: tokenMaker,
	}

	// Set up routes and run the server
	server.setRoutes()
	server.router.Run(address)

	return server, nil
}

// @ Function will configure API routes
func (server *Server) setRoutes() {
	router := gin.Default()

	// Group routes with authentication middleware
	auth := router.Group("/").Use(authMiddleware(*server.tokenMaker))
	auth.DELETE("/delete/:id", server.deleteUser) // middleware is applied only to this route in this example

	// these endpoints do not require auth
	router.POST("/create", server.createUser)
	router.POST("/login", server.login)

	// set the server router
	server.router = router
}
