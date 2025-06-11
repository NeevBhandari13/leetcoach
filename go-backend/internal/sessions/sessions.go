package sessions

import (
	"math/rand"
	"sync"

	"github.com/google/uuid"
	"github.com/neevbhandari13/leetcoach/internal/interview"
	"github.com/neevbhandari13/leetcoach/internal/models"
)

var (
	// list of strings of problem texts
	problemTexts []string = []string{"Problem text 1", "Problem text 2", "Problem text 3"}
	// map to store sessions
	sessions = make(map[string]*models.Session)
	// mutex for sessions map
	sessionsMutex sync.Mutex
)

// generateSessionID returns a new random session ID string.
func generateSessionID() string {
	return uuid.NewString()
}

func getProblemText() string {
	randomIndex := rand.Intn(len(problemTexts))
	var problemText string = problemTexts[randomIndex]
	return problemText
}

func CreateSession() *models.Session {
	session := &models.Session{
		SessionID:   generateSessionID(), // helper function
		State:       interview.IntroState,
		ChatHistory: []models.Message{},
		ProblemText: getProblemText(),
	}

	sessionsMutex.Lock()
	sessions[session.SessionID] = session
	sessionsMutex.Unlock()

	return session
}
