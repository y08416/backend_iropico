package repositories

import (
	"context"
	"fmt"

	"backend_iropico/internal/db"
	"backend_iropico/internal/models"
)

// UUID で upsert（同じuuidならname更新して返す）
func UpsertUserByUUID(ctx context.Context, name, uuid string) (models.User, error) {
	q := `
		INSERT INTO users (name, uuid)
		VALUES ($1, $2)
		ON CONFLICT (uuid) DO UPDATE SET name = EXCLUDED.name
		RETURNING id, name, uuid, created_at;
	`
	var u models.User
	err := db.DB.QueryRowContext(ctx, q, name, uuid).
		Scan(&u.ID, &u.Name, &u.UUID, &u.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to upsert user: %w", err)
	}
	return u, nil
}

func GetUserByUUID(ctx context.Context, uuid string) (models.User, error) {
	q := `SELECT id, name, uuid, created_at FROM users WHERE uuid=$1;`
	var u models.User
	err := db.DB.QueryRowContext(ctx, q, uuid).
		Scan(&u.ID, &u.Name, &u.UUID, &u.CreatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user by UUID: %w", err)
	}
	return u, nil
}
