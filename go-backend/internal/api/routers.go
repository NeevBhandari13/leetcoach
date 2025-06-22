package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/neevbhandari13/leetcoach/internal/ai"
	"github.com/neevbhandari13/leetcoach/internal/interview"
)

func SetupMiddleware(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},            // allows requests from localhost:3000, frontend server
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},           // sets methods that can be sent
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"}, // sets headers that can be sent
		AllowCredentials: true,                                         // allows cookies and auth headers in requests
		MaxAge:           12 * time.Hour,                               // cache preflight requests for 12 hours which is the check for permissions
	}))
	// add more middleware here if needed
}

// gin.Engine is the main router for Gin, which handles HTTP requests
func SetupRoutes(router *gin.Engine, gptClient *ai.GPTClient, sessionStore *interview.SessionStore) {
	// defines a post request to the /start-interview endpoint which will call the startInterviewHandler function
	router.GET("/test", testHandler)
	router.GET("/start-interview", startInterviewHandler(gptClient, sessionStore))
	router.POST("/continue-interview", continueInterviewHandler(gptClient, sessionStore))
}
