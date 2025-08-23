package repositories

import (
	"context"

	"backend_iropico/internal/db"
	"backend_iropico/internal/models"
)

type RankingRow struct {
	Uuid  string  `json:"uuid"`
	Name  string  `json:"name"`
	Score float32 `json:"score"`
}

// 1人1ラウンド1件の提出を upsert（写真はnil可）
func UpsertSubmission(
	ctx context.Context,
	roomID int64, round int, uuid string,
	h int, s, v float32,
	score float32,
	photoMime, photoURL *string,
	resultJSON []byte,
) error {
	q := `
				INSERT INTO history
					(room_id, round_no, uuid, target_h, target_s, target_v, score, result_json, photo_mime, photo_url)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
				ON CONFLICT (room_id, round_no, uuid)
				DO UPDATE SET
					target_h=EXCLUDED.target_h,
					target_s=EXCLUDED.target_s,
					target_v=EXCLUDED.target_v,
					score=EXCLUDED.score,
					result_json=EXCLUDED.result_json,
					photo_mime=EXCLUDED.photo_mime,
					photo_url=EXCLUDED.photo_url,
					created_at=NOW();
		`
	var resultJSONArg interface{}
	if len(resultJSON) == 0 {
		resultJSONArg = nil
	} else {
		resultJSONArg = resultJSON
	}
	_, err := db.DB.ExecContext(ctx, q,
		roomID, round, uuid,
		h, s, v,
		score,
		resultJSONArg,
		photoMime,
		photoURL,
	)
	return err
}

// ラウンド別ランキング
func GetRoundRanking(ctx context.Context, roomID int64, round int) ([]RankingRow, error) {
	q := `
		SELECT h.uuid, u.name, h.score
		FROM history h
		JOIN users u ON u.uuid = h.uuid
		WHERE h.room_id = $1 AND h.round_no = $2
		ORDER BY h.score DESC, h.created_at ASC;
	`
	rows, err := db.DB.QueryContext(ctx, q, roomID, round)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []RankingRow
	for rows.Next() {
		var r RankingRow
		if err := rows.Scan(&r.Uuid, &r.Name, &r.Score); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// 合計ランキング（部屋内）
func GetTotalRanking(ctx context.Context, roomID int64) ([]RankingRow, error) {
	q := `
		SELECT h.uuid, u.name, SUM(h.score) AS total_score
		FROM history h
		JOIN users u ON u.uuid = h.uuid
		WHERE h.room_id = $1
		GROUP BY h.uuid, u.name
		ORDER BY total_score DESC;
	`
	rows, err := db.DB.QueryContext(ctx, q, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []RankingRow
	for rows.Next() {
		var r RankingRow
		if err := rows.Scan(&r.Uuid, &r.Name, &r.Score); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// ユーザーの履歴（最新順）
func GetUserHistory(ctx context.Context, uuid string) ([]models.History, error) {
	q := `
		SELECT id, room_id, round_no, uuid,
			   target_h, target_s, target_v,
			   score, result_json, photo_mime, photo_url, created_at
		FROM history
		WHERE uuid=$1
		ORDER BY created_at DESC;
	`
	rows, err := db.DB.QueryContext(ctx, q, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.History
	for rows.Next() {
		var h models.History
		if err := rows.Scan(
			&h.ID, &h.RoomID, &h.RoundNo, &h.UserID,
			&h.TargetH, &h.TargetS, &h.TargetV,
			&h.Score, &h.ResultJSON, &h.PhotoMime, &h.PhotoURL, &h.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, rows.Err()
}
