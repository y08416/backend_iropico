package models

import "time"

type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	UUID      string    `db:"uuid"` // フロント生成のUUID（DB型: uuid）
	CreatedAt time.Time `db:"created_at"`
}
