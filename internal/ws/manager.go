// internal/ws/manager.go
package ws

import "sync"

type Manager struct {
	mu   sync.RWMutex
	hubs map[string]*Hub // key: room code
}

func NewManager() *Manager { return &Manager{hubs: make(map[string]*Hub)} }

func (m *Manager) GetHub(code string) *Hub {
	m.mu.RLock()
	h := m.hubs[code]
	m.mu.RUnlock()
	if h != nil {
		return h
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.hubs[code] == nil {
		m.hubs[code] = NewHub()
	}
	return m.hubs[code]
}

func (m *Manager) DeleteHub(code string) {
	m.mu.Lock()
	delete(m.hubs, code)
	m.mu.Unlock()
}
