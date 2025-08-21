// internal/ws/handler.go
package ws

import (
	"github.com/gofiber/websocket/v2"
)

func ServeWS(m *Manager) func(*websocket.Conn) {
	return func(conn *websocket.Conn) {
		// /ws/rooms/:code
		roomCode := conn.Params("code")
		// 任意: ?user_id=&uuid=
		_ = conn.Query("user_id")
		_ = conn.Query("uuid")

		hub := m.GetHub(roomCode)

		// クライアント登録
		hub.Register(conn)
		defer hub.Unregister(conn)

		// サーバ主導: 受信は必要最低限だけ読む（切断検知のため）
		for {
			var in map[string]any
			if err := conn.ReadJSON(&in); err != nil {
				// ここにログを入れてもOK（log.Printfなど）
				return // クライアント切断
			}
			// 例）クライアントからのpingに応答したい場合:
			// if in["type"] == "ping" {
			// 	hub.Emit("pong", nil)
			// }
		}
	}
}
