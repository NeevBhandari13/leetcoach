package api

import (
	"github.com/NeevBhandari13/leetcoach/internal/chat"
	"github.com/NeevBhandari13/leetcoach/internal/llm"
	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	System   string        `json:"system"`
	Messages []llm.Message `json:"messages"`
}

type ChatResponse struct {
	Message string `json:"message"`
}

func ChatHandler(chatService *chat.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
