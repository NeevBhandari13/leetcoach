package main

import (
	"github.com/gin-gonic/gin"
	"github.com/neevbhandari13/leetcoach/internal/api"
)

func main() {
	router := gin.Default() // starts default Gin router with logging and recovery
	api.SetupRoutes(router) // sets up the routes for the API
	router.Run(":8080")     // starts the server on port 8080
}
