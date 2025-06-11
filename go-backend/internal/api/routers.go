package api

import (
	"github.com/gin-gonic/gin"
)

// gin.Engine is the main router for Gin, which handles HTTP requests
func SetupRoutes(router *gin.Engine) {
	// defines a post request to the /start-interview endpoint which will call the startInterviewHandler function
	router.GET("/test", testHandler)
	router.POST("/start-interview", startInterviewHandler)
}
