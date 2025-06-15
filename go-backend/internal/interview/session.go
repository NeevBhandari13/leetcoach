package interview

import (
	"fmt"
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

func GetSession(sessionID string) (*models.Session, error) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	session, ok := sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session with ID %s not found", sessionID)
	}
	return session, nil
}

func GetState(sessionID string) (models.State, error) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	session, ok := sessions[sessionID]
	if !ok {
		return models.NilState, fmt.Errorf("session with ID %s not found", sessionID)
	}
	return session.State, nil
}

func GetProblemText(sessionID string) (string, error) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	session, ok := sessions[sessionID]
	if !ok {
		return "", fmt.Errorf("session with ID %s not found", sessionID)
	}
	return session.ProblemText, nil
}
