package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"backend_iropico/internal/db"
	"backend_iropico/internal/models"
)

func CreateRoom(ctx context.Context, uid string, code string) (models.Room, error) {
	var hostUuid string
	err := db.DB.QueryRowContext(ctx, "SELECT uuid FROM users WHERE uuid = $1", uid).Scan(&hostUuid)
	if err != nil {
		return models.Room{}, fmt.Errorf("failed to get user UUID: %w", err)
	}

	q := `
    INSERT INTO rooms (code, host_uuid)
    VALUES ($1, $2)
    RETURNING id, code, host_uuid, status, current_round, created_at, started_at, ended_at;
  `
	var r models.Room
	err = db.DB.QueryRowContext(ctx, q, code, hostUuid).
		Scan(&r.ID, &r.Code, &r.HostUuid, &r.Status, &r.CurrentRound, &r.CreatedAt, &r.StartedAt, &r.EndedAt)
	if err != nil {
		return models.Room{}, fmt.Errorf("failed to create room: %w", err)
	}

	return r, nil
}

func GetRoomByCode(ctx context.Context, code string) (models.Room, error) {
	q := `
		SELECT id, code, host_uuid, status, current_round, created_at, started_at, ended_at
		FROM rooms WHERE code=$1;
	`
	var r models.Room
	err := db.DB.QueryRowContext(ctx, q, code).
		Scan(&r.ID, &r.Code, &r.HostUuid, &r.Status, &r.CurrentRound, &r.CreatedAt, &r.StartedAt, &r.EndedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return r, nil
		}
		return r, err
	}
	return r, nil
}

func GetUserIDByUUID(ctx context.Context, uuid string) (int64, error) {
	var userID int64
	q := `
		SELECT id
		FROM users
		WHERE uuid = $1;
	`
	err := db.DB.QueryRowContext(ctx, q, uuid).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func UpdateRoomStatusAndRound(ctx context.Context, roomID int64, status string, round int) error {
	q := `UPDATE rooms SET status=$1, current_round=$2 WHERE id=$3;`
	res, err := db.DB.ExecContext(ctx, q, status, round, roomID)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("room not found")
	}
	return nil
}

func SetRoomStarted(ctx context.Context, roomID int64) error {
	q := `UPDATE rooms SET status='started', current_round=1, started_at=NOW() WHERE id=$1;`
	_, err := db.DB.ExecContext(ctx, q, roomID)
	return err
}

func IncrementRound(ctx context.Context, roomID int64) (int, error) {
	q := `UPDATE rooms SET current_round = current_round + 1 WHERE id=$1 RETURNING current_round;`
	var round int
	err := db.DB.QueryRowContext(ctx, q, roomID).Scan(&round)
	return round, err
}

func SetRoomEnded(ctx context.Context, roomID int64) error {
	q := `UPDATE rooms SET status='ended', ended_at=NOW() WHERE id=$1;`
	_, err := db.DB.ExecContext(ctx, q, roomID)
	return err
}
