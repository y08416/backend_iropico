package models

import "time"

type Player struct {
	ID         int64     `db:"id"`
	RoomID     int64     `db:"room_id"`
	UserID     int64     `db:"user_id"`
	HasCleared bool      `db:"has_cleared"` // ラウンド毎にリセット
	JoinedAt   time.Time `db:"joined_at"`
}
