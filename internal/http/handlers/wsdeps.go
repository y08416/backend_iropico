package handlers

import "backend_iropico/internal/ws"

// パッケージ内で使うWSマネージャ（MVP用の簡易DI）
var wsMgr *ws.Manager

// ルーター側から一度だけ注入する
func SetWSManager(m *ws.Manager) {
	wsMgr = m
}
