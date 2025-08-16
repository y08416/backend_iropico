package ws

// Manager はWebSocketの接続管理用のダミー構造体
type Manager struct{}

// NewManager はManagerの新しいインスタンスを返す
func NewManager() *Manager {
	return &Manager{}
}
