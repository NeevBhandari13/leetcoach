package api

import (
	"net/http"

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
		// variable to hold the request body
		var req ChatRequest

		// tries to map request into ChatRequest
		// if the request body doesnt match the definition, we
		// throw an error
		err := c.ShouldBindJSON(&req)
		if err != nil {
			// gin.H creates a simple json
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		resp, err := chatService.Reply(c, req.System, req.Messages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, ChatResponse{
			Message: resp,
		})
		return
	}
}
