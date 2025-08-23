package handlers

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"backend_iropico/internal/repositories"
)

func Submit(c *fiber.Ctx) error {
	code := c.Params("code")

	// multipart: user_id, round_no, score, h, s, v, photo(file)
	uuid := c.FormValue("uuid")
	roundNo, _ := strconv.Atoi(c.FormValue("round_no"))
	scoreF, _ := strconv.ParseFloat(c.FormValue("score"), 32)
	h, _ := strconv.Atoi(c.FormValue("h"))
	sF, _ := strconv.ParseFloat(c.FormValue("s"), 32)
	vF, _ := strconv.ParseFloat(c.FormValue("v"), 32)

	if uuid == "" || roundNo == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "msg": "uuid and round_no required"})
	}

	// 画像（任意）
	var photoMime *string
	var photoURL *string // 今は未使用
	if fileHeader, err := c.FormFile("photo"); err == nil && fileHeader != nil {
		f, err := fileHeader.Open()
		if err == nil {
			defer f.Close()
			// photoBytesは不要
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
		room.ID, roundNo, uuid,
		h, float32(sF), float32(vF),
		float32(scoreF),
		photoMime, photoURL,
		nil, // resultJSON (任意)
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": err.Error()})
	}

	// そのユーザーを「クリア」扱いに
	if err := repositories.MarkCleared(ctx, room.ID, uuid); err != nil {
		log.Printf("[WARN] MarkCleared failed: %v", err)
	}

	// ラウンドの最新ランキングを返す
	rank, err := repositories.GetRoundRanking(ctx, room.ID, roundNo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": err.Error()})
	}

	// ★ WS通知：提出/ランキング更新
	if wsMgr != nil {
		wsMgr.GetHub(code).Emit("submission", fiber.Map{
			"round":   roundNo,
			"uuid":    uuid,
			"score":   float32(scoreF),
			"ranking": rank,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"type":    "submit",
		"ok":      true,
		"round":   roundNo,
		"uuid":    uuid,
		"score":   float32(scoreF),
		"ranking": rank,
	})
}
