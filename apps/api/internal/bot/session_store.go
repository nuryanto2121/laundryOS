package bot

import (
	"sync"
	"time"
)

// SessionStore defines the contract for persisting conversation states.
type SessionStore interface {
	Get(phone string) (*SessionData, bool)
	Set(phone string, data *SessionData)
	Delete(phone string)
}

// MemorySessionStore is a thread-safe in-memory implementation of SessionStore.
type MemorySessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*SessionData
	ttl      time.Duration
}

// NewMemorySessionStore creates a new MemorySessionStore with the given TTL.
func NewMemorySessionStore(ttl time.Duration) *MemorySessionStore {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	store := &MemorySessionStore{
		sessions: make(map[string]*SessionData),
		ttl:      ttl,
	}

	// Background routine for periodic cleanup of expired sessions
	go store.cleanupLoop(5 * time.Minute)

	return store
}

// Get retrieves the session data for a phone number if it exists and has not expired.
func (s *MemorySessionStore) Get(phone string) (*SessionData, bool) {
	s.mu.RLock()
	session, exists := s.sessions[phone]
	s.mu.RUnlock()

	if !exists {
		return nil, false
	}

	if session.IsExpired(s.ttl) {
		s.Delete(phone)
		return nil, false
	}

	return session, true
}

// Set stores or updates the session data for a phone number.
func (s *MemorySessionStore) Set(phone string, data *SessionData) {
	if data == nil {
		s.Delete(phone)
		return
	}
	data.Phone = phone
	data.LastActive = time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[phone] = data
}

// Delete removes the session for the given phone number.
func (s *MemorySessionStore) Delete(phone string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, phone)
}

func (s *MemorySessionStore) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		s.mu.Lock()
		for phone, session := range s.sessions {
			if session.IsExpired(s.ttl) {
				delete(s.sessions, phone)
			}
		}
		s.mu.Unlock()
	}
}
