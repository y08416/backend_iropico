package repositories

import (
	"context"

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
	return u, err
}

func GetUserByID(ctx context.Context, id int64) (models.User, error) {
	q := `SELECT id, name, uuid, created_at FROM users WHERE id=$1;`
	var u models.User
	err := db.DB.QueryRowContext(ctx, q, id).
		Scan(&u.ID, &u.Name, &u.UUID, &u.CreatedAt)
	return u, err
}
