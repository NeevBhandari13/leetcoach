package api

import (
	"github.com/gin-gonic/gin"
	"github.com/neevbhandari13/leetcoach/internal/ai"
	"github.com/neevbhandari13/leetcoach/internal/interview"
)

// gin.Engine is the main router for Gin, which handles HTTP requests
func SetupRoutes(router *gin.Engine, gptClient *ai.GPTClient, sessionStore *interview.SessionStore) {
	// defines a post request to the /start-interview endpoint which will call the startInterviewHandler function
	router.GET("/test", testHandler)
	router.GET("/start-interview", startInterviewHandler(gptClient, sessionStore))
	router.POST("/continue-interview", continueInterviewHandler(gptClient, sessionStore))
}
