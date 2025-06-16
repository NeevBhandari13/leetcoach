package interview

import (
	"fmt"
	"sync"

	"github.com/neevbhandari13/leetcoach/internal/models"
	"github.com/neevbhandari13/leetcoach/internal/utils"
	"github.com/neevbhandari13/leetcoach/pkg/problems"
)

type SessionStore struct {
	// map to store sessions
	sessions map[string]*models.Session
	// mutex for sessions map
	mu sync.Mutex
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*models.Session),
	}
}

func (s *SessionStore) CreateSession() *models.Session {
	session := &models.Session{
		SessionID:   utils.GenerateSessionID(), // helper function
		State:       models.IntroState,
		ChatHistory: []models.Message{},
		ProblemText: problems.GetProblemText(),
	}

	s.mu.Lock()
	s.sessions[session.SessionID] = session
	s.mu.Unlock()

	return session
}

func (s *SessionStore) GetSession(sessionID string) (*models.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session with ID %s not found", sessionID)
	}
	return session, nil
}

func (s *SessionStore) GetState(sessionID string) (models.State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return models.NilState, fmt.Errorf("session with ID %s not found", sessionID)
	}
	return session.State, nil
}

func (s *SessionStore) SetState(sessionID string, state models.State) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return fmt.Errorf("session with ID %s not found", sessionID)
	}
	session.State = state
	return nil
}

func (s *SessionStore) GetProblemText(sessionID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return "", fmt.Errorf("session with ID %s not found", sessionID)
	}
	return session.ProblemText, nil
}
