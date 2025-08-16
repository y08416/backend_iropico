package handlers

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"backend_iropico/internal/repositories"
)

func Submit(c *fiber.Ctx) error {
	code := c.Params("code")

	// multipart: user_id, round_no, score, h, s, v, photo(file)
	userID, _ := strconv.ParseInt(c.FormValue("user_id"), 10, 64)
	roundNo, _ := strconv.Atoi(c.FormValue("round_no"))
	scoreF, _ := strconv.ParseFloat(c.FormValue("score"), 32)
	h, _ := strconv.Atoi(c.FormValue("h"))
	sF, _ := strconv.ParseFloat(c.FormValue("s"), 32)
	vF, _ := strconv.ParseFloat(c.FormValue("v"), 32)

	if userID == 0 || roundNo == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "msg": "user_id and round_no required"})
	}

	// 画像（任意）
	var photoBytes []byte
	var photoMime *string
	var photoURL *string // 今は未使用
	if fileHeader, err := c.FormFile("photo"); err == nil && fileHeader != nil {
		f, err := fileHeader.Open()
		if err == nil {
			defer f.Close()
			b, _ := io.ReadAll(f)
			photoBytes = b
			m := fileHeader.Header.Get("Content-Type")
			photoMime = &m
		}
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// ルーム取得
	room, err := repositories.GetRoomByCode(ctx, code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"ok": false, "msg": "room not found"})
	}

	// 保存（upsert）
	err = repositories.UpsertSubmission(
		ctx,
		room.ID, roundNo, userID,
		h, float32(sF), float32(vF),
		float32(scoreF),
		photoBytes, photoMime, photoURL,
		nil, // resultJSON (任意)
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": err.Error()})
	}

	// そのユーザーを「クリア」扱いに
	if err := repositories.MarkCleared(ctx, room.ID, userID); err != nil {
		// 参加してない等で失敗しても致命的ではないので警告にとどめる運用でもOK
	}

	// ラウンドの最新ランキングを返す
	rank, err := repositories.GetRoundRanking(ctx, room.ID, roundNo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok":      true,
		"round":   roundNo,
		"ranking": rank,
	})
}
