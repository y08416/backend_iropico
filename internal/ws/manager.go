package ws

import "sync"

type Manager struct {
	mu   sync.Mutex
	hubs map[string]*Hub
}

func NewManager() *Manager {
	return &Manager{hubs: make(map[string]*Hub)}
}

func (m *Manager) GetHub(code string) *Hub {
	m.mu.Lock()
	defer m.mu.Unlock()

	if h, ok := m.hubs[code]; ok {
		return h
	}
	h := NewHub()
	m.hubs[code] = h
	return h
}
