// Package voice proporciona el middleware reutilizable para el agente de voz inteligente.
package voice

import (
	"fmt"
	"sync"
)

// SessionStore interfaz para manejar sesiones
type SessionStore interface {
	Get(id string) (*Session, error)
	Save(session *Session) error
	Delete(id string) error
}

// MemoryStore implementación en memoria de SessionStore
type MemoryStore struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

// NewMemoryStore crea un nuevo MemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		sessions: make(map[string]*Session),
	}
}

// Get obtiene una sesión por ID
func (m *MemoryStore) Get(id string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if s, ok := m.sessions[id]; ok {
		return s, nil
	}

	// Crear nueva sesión si no existe
	session := &Session{
		ID:            id,
		State:         StateIdle,
		Context:       make(map[string]string),
		TemporaryData: make(map[string]interface{}),
	}

	m.sessions[id] = session
	return session, nil
}

// Save guarda una sesión
func (m *MemoryStore) Save(session *Session) error {
	if session == nil {
		return fmt.Errorf("sesión no puede ser nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[session.ID] = session
	return nil
}

// Delete elimina una sesión
func (m *MemoryStore) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, id)
	return nil
}
