package models

import "time"

// status は DB では ENUM (waiting|started|ended)
type Room struct {
	ID           int64      `db:"id"`
	Code         string     `db:"code"` // 参加用コード（UNIQUE）
	HostUserID   int64      `db:"host_user_id"`
	Status       string     `db:"status"`        // "waiting" / "started" / "ended"
	CurrentRound int        `db:"current_round"` // 0 開始、開始後に 1 になる
	CreatedAt    time.Time  `db:"created_at"`
	StartedAt    *time.Time `db:"started_at"` // NULL 可
	EndedAt      *time.Time `db:"ended_at"`   // NULL 可
}
