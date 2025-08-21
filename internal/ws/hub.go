package ws

import "sync"

type client interface {
	WriteJSON(v interface{}) error
	Close() error
}

type Hub struct {
	mu      sync.RWMutex
	clients map[client]struct{}
}

func NewHub() *Hub {
	return &Hub{clients: make(map[client]struct{})}
}

func (h *Hub) Register(c client) {
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()
}

func (h *Hub) Unregister(c client) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
	_ = c.Close()
}

func (h *Hub) BroadcastJSON(v interface{}) {
	h.mu.RLock()
	for c := range h.clients {
		_ = c.WriteJSON(v)
	}
	h.mu.RUnlock()
}

func (h *Hub) Emit(event string, payload any) {
	h.BroadcastJSON(map[string]any{
		"type":    event,
		"payload": payload,
	})
}
