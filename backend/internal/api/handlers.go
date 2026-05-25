package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/NeevBhandari13/leetcoach/internal/chat"
	"github.com/NeevBhandari13/leetcoach/internal/llm"
	"github.com/NeevBhandari13/leetcoach/internal/prompts"
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

// Store is the subset of session.SessionStore the handlers need. Using an
// interface here lets tests inject a mock without a real database.
type Store interface {
	CreateSession(ctx context.Context, sessionID string) (*session.Session, error)
	UpdateChatHistory(ctx context.Context, sessionID string, message llm.Message) error
	GetSession(ctx context.Context, sessionID string) (*session.Session, error)
	Reply(ctx context.Context, sessionID, system, userMessage string) (string, error)
	SetState(ctx context.Context, sessionID string, state session.State) error
}

// llmReply is the JSON structure the LLM is instructed to return on every turn.
type llmReply struct {
	Reply        string        `json:"reply"`
	CurrentState session.State `json:"current_state"`
}

type replyRequest struct {
	Message string `json:"message"`
}

type replyResponse struct {
	Message string `json:"message"`
}

func ReplyHandler(store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")

		var req replyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Load the session so we know the current state and which problem is
		// being discussed. Both are needed to build the correct system prompt.
		sess, err := store.GetSession(c, sessionID)
		if err != nil {
			if errors.Is(err, session.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Build the system prompt from the session's current state and problem
		// text. This keeps the LLM's instructions entirely server-side.
		systemPrompt := prompts.GetSystemPrompt(sess.State, sess.ProblemText)

		// Persist the user message, load full history, call the LLM, persist
		// the assistant reply — all handled inside store.Reply.
		raw, err := store.Reply(c, sessionID, systemPrompt, req.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// The LLM is instructed to return JSON with 'reply' and 'current_state'.
		// We parse it here so we can advance the state machine and return only
		// the human-readable reply to the client.
		var llmRes llmReply
		if err := json.Unmarshal([]byte(raw), &llmRes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse LLM response"})
			return
		}

		// Advance the session state based on what the LLM decided.
		if err := store.SetState(c, sessionID, llmRes.CurrentState); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, replyResponse{Message: llmRes.Reply})
	}
}

func StartInterviewHandler(store Store) gin.HandlerFunc {
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
