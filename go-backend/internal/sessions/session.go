package sessions

import (
	"sync"

	"github.com/google/uuid"
	"github.com/neevbhandari13/leetcoach/internal/models"
	"github.com/neevbhandari13/leetcoach/pkg/problems"
)

var (
	// map to store sessions
	sessions = make(map[string]*models.Session)
	// mutex for sessions map
	sessionsMutex sync.Mutex
)

// generateSessionID returns a new random session ID string.
func generateSessionID() string {
	return uuid.NewString()
}

func CreateSession() *models.Session {
	session := &models.Session{
		SessionID:   generateSessionID(), // helper function
		State:       models.IntroState,
		ChatHistory: []models.Message{},
		ProblemText: problems.GetProblemText(),
	}

	sessionsMutex.Lock()
	sessions[session.SessionID] = session
	sessionsMutex.Unlock()

	return session
}

// Gets chat history from sessionID
func GetChatHistory(sessionID string) []models.Message {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	return sessions[sessionID].ChatHistory
}

func UpdateChatHistory(sessionID string, message models.Message) {
	// lock the sessions map
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	// get chat history from session and append the new message to it
	sessions[sessionID].ChatHistory = append(sessions[sessionID].ChatHistory, message)
}
