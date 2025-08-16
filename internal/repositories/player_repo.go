package repositories

import (
	"context"
	"errors"

	"backend_iropico/internal/db"
	"backend_iropico/internal/models"
)

func JoinRoom(ctx context.Context, roomID, userID int64) error {
	// UNIQUE (room_id, user_id) なので二重参加はDB側で弾かれる
	q := `INSERT INTO players (room_id, user_id) VALUES ($1, $2)
	      ON CONFLICT (room_id, user_id) DO NOTHING;`
	_, err := db.DB.ExecContext(ctx, q, roomID, userID)
	return err
}

func GetPlayersByRoom(ctx context.Context, roomID int64) ([]models.Player, error) {
	q := `SELECT id, room_id, user_id, has_cleared, joined_at
	      FROM players WHERE room_id=$1 ORDER BY joined_at ASC;`
	rows, err := db.DB.QueryContext(ctx, q, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.Player{}
	for rows.Next() {
		var p models.Player
		if err := rows.Scan(&p.ID, &p.RoomID, &p.UserID, &p.HasCleared, &p.JoinedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func ResetHasCleared(ctx context.Context, roomID int64) error {
	q := `UPDATE players SET has_cleared=false WHERE room_id=$1;`
	_, err := db.DB.ExecContext(ctx, q, roomID)
	return err
}

func MarkCleared(ctx context.Context, roomID, userID int64) error {
	q := `UPDATE players SET has_cleared=true WHERE room_id=$1 AND user_id=$2;`
	res, err := db.DB.ExecContext(ctx, q, roomID, userID)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("player not found")
	}
	return nil
}
