package models

import (
	"time"
)

type History struct {
	ID      int64 `db:"id"`
	RoomID  int64 `db:"room_id"`
	RoundNo int   `db:"round_no"`
	UserID  int64 `db:"user_id"`

	TargetH int     `db:"target_h"` // 0..360
	TargetS float32 `db:"target_s"` // 0..1
	TargetV float32 `db:"target_v"` // 0..1

	Score      float32 `db:"score"`       // 0..1
	ResultJSON []byte  `db:"result_json"` // JSONB は []byte で受ける or stringでも可

	PhotoBytes []byte  `db:"photo_bytes"` // MVPはDB保存
	PhotoMime  *string `db:"photo_mime"`
	PhotoURL   *string `db:"photo_url"`

	CreatedAt time.Time `db:"created_at"`
}
