package ws

import (
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

type safeConn struct {
	c  *websocket.Conn
	mu sync.Mutex
}

func newSafeConn(c *websocket.Conn) *safeConn {
	return &safeConn{c: c}
}

func (s *safeConn) WriteJSON(v interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.c.WriteJSON(v)
}

func (s *safeConn) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.c.Close()
}

func (s *safeConn) ping() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.c.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(10*time.Second))
}
