package utils

import "github.com/google/uuid"

// generateSessionID returns a new random session ID string.
func GenerateSessionID() string {
	return uuid.NewString()
}
