// internal/ws/hub.go
package ws

import (
	"sync"
)

type client interface {
	WriteJSON(v interface{}) error
	Close() error
}

type Hub struct {
	mu      sync.RWMutex
	clients map[client]struct{}
	closed  bool
}

func NewHub() *Hub { return &Hub{clients: make(map[client]struct{})} }

func (h *Hub) Register(c client) {
	h.mu.Lock()
	if !h.closed {
		h.clients[c] = struct{}{}
	}
	h.mu.Unlock()
}

func (h *Hub) Unregister(c client) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
	_ = c.Close()
}

func (h *Hub) BroadcastJSON(v interface{}) {
	// 失敗しても他クライアントに影響出さない
	h.mu.RLock()
	for c := range h.clients {
		_ = c.WriteJSON(v)
	}
	h.mu.RUnlock()
}

// 便利関数：type+payload で送る
func (h *Hub) Emit(event string, payload any) {
	msg := map[string]any{"type": event, "payload": payload}
	h.BroadcastJSON(msg)
}
