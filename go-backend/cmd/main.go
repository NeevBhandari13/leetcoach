package main

import (
	"github.com/gin-gonic/gin"
	"github.com/neevbhandari13/leetcoach/internal/ai"
	"github.com/neevbhandari13/leetcoach/internal/api"
	"github.com/neevbhandari13/leetcoach/internal/interview"
)

func main() {
	// initialises new gpt client
	gptClient := ai.NewGPTClient()

	// initialises new session storage
	sessionStore := interview.NewSessionStore() // initialises new session store correctly with empty map

	router := gin.Default()                          // starts default Gin router with logging and recovery
	api.SetupMiddleware(router)                      // sets up middleware for gin router
	api.SetupRoutes(router, gptClient, sessionStore) // sets up the routes for the API
	router.Run(":8080")                              // starts the server on port 8080
}
