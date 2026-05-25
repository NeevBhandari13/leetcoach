package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(store Store) *gin.Engine {
	router := gin.Default()
	router.HandleMethodNotAllowed = true
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		MaxAge:           12 * time.Hour,
	}))
	SetupRoutes(router, store)
	return router
}

func SetupRoutes(router *gin.Engine, store Store) {
	router.POST("/start", StartInterviewHandler(store))
	router.GET("/sessions/:id", GetSessionHandler(store))
	router.POST("/sessions/:id/reply", ReplyHandler(store))
}
