package memory

import (
	"github.com/neevbhandari13/leetcoach/internal/interview"
	"sync"
)

var (
	sessions      = make(map[string]*interview.Session)
	sessionsMutex sync.Mutex
)
