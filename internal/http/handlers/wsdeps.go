package handlers

import "backend_iropico/internal/ws"

var wsMgr *ws.Manager

func SetWSManager(m *ws.Manager) {
	wsMgr = m
}
