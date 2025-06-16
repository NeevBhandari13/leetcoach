package api

import (
	"github.com/gin-gonic/gin"
	"github.com/neevbhandari13/leetcoach/internal/ai"
)

// gin.Engine is the main router for Gin, which handles HTTP requests
func SetupRoutes(router *gin.Engine, gptClient *ai.GPTClient) {
	// defines a post request to the /start-interview endpoint which will call the startInterviewHandler function
	router.GET("/test", testHandler)
	router.POST("/start-interview", startInterviewHandler(gptClient))
	router.POST("/continue-interview", continueInterviewHandler(gptClient))
}
