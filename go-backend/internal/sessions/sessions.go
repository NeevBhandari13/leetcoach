package sessions

import (
	"math/rand"
	"sync"

	"github.com/google/uuid"
	"github.com/neevbhandari13/leetcoach/internal/interview"
	"github.com/neevbhandari13/leetcoach/pkg/problems"
)

var (
	// map to store sessions
	sessions = make(map[string]*interview.Session)
	// mutex for sessions map
	sessionsMutex sync.Mutex
)

// generateSessionID returns a new random session ID string.
func generateSessionID() string {
	return uuid.NewString()
}

func getProblemText() string {
	randomIndex := rand.Intn(len(problems.ProblemTexts))
	var problemText string = problems.problemTexts[randomIndex]
	return problemText
}

func CreateSession() *interview.Session {
	session := &interview.Session{
		SessionID:   generateSessionID(), // helper function
		State:       interview.IntroState,
		ChatHistory: []interview.Message{},
		ProblemText: getProblemText(),
	}

	sessionsMutex.Lock()
	sessions[session.SessionID] = session
	sessionsMutex.Unlock()

	return session
}
