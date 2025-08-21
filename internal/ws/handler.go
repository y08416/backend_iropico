package ws

import (
	"time"

	"github.com/gofiber/websocket/v2"
)

func ServeWS(m *Manager) func(*websocket.Conn) {
	return func(conn *websocket.Conn) {
		roomCode := conn.Params("code")
		_ = conn.Query("user_id")
		_ = conn.Query("uuid")

		conn.SetReadLimit(1 << 20)
		_ = conn.SetReadDeadline(time.Now().Add(70 * time.Second))
		conn.SetPongHandler(func(string) error {
			return conn.SetReadDeadline(time.Now().Add(70 * time.Second))
		})

		sc := newSafeConn(conn)
		hub := m.GetHub(roomCode)
		hub.Register(sc)
		defer hub.Unregister(sc)

		stop := make(chan struct{})
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					_ = sc.ping()
				case <-stop:
					return
				}
			}
		}()

		for {
			var in map[string]any
			if err := conn.ReadJSON(&in); err != nil {
				close(stop)
				return
			}
		}
	}
}
