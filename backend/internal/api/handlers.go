package api

import (
	"net/http"

	"github.com/NeevBhandari13/leetcoach/internal/chat"
	"github.com/NeevBhandari13/leetcoach/internal/llm"
	"github.com/NeevBhandari13/leetcoach/internal/session"
	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	System   string        `json:"system"`
	Messages []llm.Message `json:"messages"`
}

type ChatResponse struct {
	Message string `json:"message"`
}

type startInterviewResponse struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

func ChatHandler(chatService *chat.ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := chatService.Reply(c, req.System, req.Messages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ChatResponse{Message: resp})
	}
}

func StartInterviewHandler(store *session.SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, err := store.CreateSession(c, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		reply := "Hello! Welcome to LeetCoach! Today we are going to be running a technical interview going over a coding problem together. Remember I'm here to help you! Are you ready to get started?"

		if err := store.UpdateChatHistory(c, sess.SessionID, llm.Message{Role: "assistant", Content: reply}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, startInterviewResponse{
			SessionID: sess.SessionID,
			Message:   reply,
		})
	}
}
