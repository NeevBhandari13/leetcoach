package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/NeevBhandari13/leetcoach/internal/llm"
	"github.com/NeevBhandari13/leetcoach/internal/prompts"
	"github.com/NeevBhandari13/leetcoach/internal/session"
	"github.com/gin-gonic/gin"
)

// Store is the subset of session.SessionStore the handlers need. Using an
// interface here lets tests inject a mock without a real database.
type Store interface {
	GetRandomProblemText(ctx context.Context) (string, error)
	CreateSession(ctx context.Context, sessionID, problemText string) (*session.Session, error)
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

type startInterviewResponse struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

type replyRequest struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type replyResponse struct {
	Message string `json:"message"`
}

type getSessionResponse struct {
	SessionID   string         `json:"session_id"`
	State       session.State  `json:"state"`
	ProblemText string         `json:"problem_text"`
	ChatHistory []llm.Message  `json:"chat_history"`
}

func StartInterviewHandler(store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pick a random problem from the DB so the LLM knows what to discuss.
		problemText, err := store.GetRandomProblemText(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		sess, err := store.CreateSession(c, "", problemText)
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

func ReplyHandler(store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")

		var req replyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sess, err := store.GetSession(c, sessionID)
		if err != nil {
			if errors.Is(err, session.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		systemPrompt := prompts.GetSystemPrompt(sess.State, sess.ProblemText, req.Code)

		raw, err := store.Reply(c, sessionID, systemPrompt, req.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var llmRes llmReply
		if err := json.Unmarshal([]byte(raw), &llmRes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse LLM response"})
			return
		}

		if err := store.SetState(c, sessionID, llmRes.CurrentState); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := store.UpdateChatHistory(c, sessionID, llm.Message{Role: "assistant", Content: llmRes.Reply}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, replyResponse{Message: llmRes.Reply})
	}
}

func GetSessionHandler(store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")

		sess, err := store.GetSession(c, sessionID)
		if err != nil {
			if errors.Is(err, session.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, getSessionResponse{
			SessionID:   sess.SessionID,
			State:       sess.State,
			ProblemText: sess.ProblemText,
			ChatHistory: sess.ChatHistory,
		})
	}
}
