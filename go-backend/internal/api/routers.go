package api

import (
	"github.com/NeevBhandari13/leetcoach/internal/chat"
	"github.com/gin-gonic/gin"
)

func NewRouter(chatService *chat.ChatService) *gin.Engine {
	router := gin.Default()
	SetupRoutes(router, chatService)
	return router
}

func SetupRoutes(router *gin.Engine, chatService *chat.ChatService) {
	router.POST("/chat", ChatHandler(chatService))
}
