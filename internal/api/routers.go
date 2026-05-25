package api

import (
	"github.com/NeevBhandari13/leetcoach/internal/chat"
	"github.com/NeevBhandari13/leetcoach/internal/session"
	"github.com/gin-gonic/gin"
)

func NewRouter(chatService *chat.ChatService, store *session.SessionStore) *gin.Engine {
	router := gin.Default()
	router.HandleMethodNotAllowed = true
	SetupRoutes(router, chatService, store)
	return router
}

func SetupRoutes(router *gin.Engine, chatService *chat.ChatService, store *session.SessionStore) {
	router.POST("/chat", ChatHandler(chatService))
	router.POST("/start", StartInterviewHandler(store))
}
