package handlers

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"backend_iropico/internal/repositories"
)

// -------- ルーム作成 --------
type createRoomReq struct {
	HostUserID int64 `json:"host_user_id"`
}

func CreateRoom(c *fiber.Ctx) error {
	var req createRoomReq
	if err := c.BodyParser(&req); err != nil || req.HostUserID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "msg": "host_user_id required"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	// 一意コードを複数回トライ（最大10回）
	for i := 0; i < 10; i++ {
		code := genCode(6)
		r, err := repositories.CreateRoom(ctx, req.HostUserID, code)
		if err == nil {
			// ※必要ならここでホストも自動参加させる
			// _ = repositories.JoinRoom(ctx, r.ID, req.HostUserID)

			return c.Status(fiber.StatusCreated).JSON(fiber.Map{
				"ok":      true,
				"room_id": r.ID,
				"code":    r.Code,
			})
		}
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": "failed to create room"})
}

// -------- ルーム参加 --------
type joinReq struct {
	UserID int64 `json:"user_id"`
}

func JoinRoom(c *fiber.Ctx) error {
	code := c.Params("code")
	var req joinReq
	if err := c.BodyParser(&req); err != nil || req.UserID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "msg": "user_id required"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	room, err := repositories.GetRoomByCode(ctx, code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"ok": false, "msg": "room not found"})
	}

	if err := repositories.JoinRoom(ctx, room.ID, req.UserID); err != nil {
		// UNIQUE制約で既参加なら実害ないので 204 でも良いが、ここではエラー返却
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": err.Error()})
	}

	// WS通知（参加者が増えた）
	if wsMgr != nil {
		hub := wsMgr.GetHub(code)
		hub.Emit("player_joined", fiber.Map{
			"userId": req.UserID,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// -------- 開始 / 次ラウンド / 終了 --------

func StartRoom(c *fiber.Ctx) error {
	code := c.Params("code")

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	room, err := repositories.GetRoomByCode(ctx, code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"ok": false, "msg": "room not found"})
	}

	// ゲーム開始: status=started, current_round=1
	if err := repositories.SetRoomStarted(ctx, room.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": err.Error()})
	}
	// 各プレイヤーのクリアフラグをリセット
	if err := repositories.ResetHasCleared(ctx, room.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": err.Error()})
	}

	// 出題色
	h, s, v := randomHSV()

	// WSで部屋全員に配信
	if wsMgr != nil {
		hub := wsMgr.GetHub(code)
		hub.Emit("game_started", fiber.Map{
			"round": 1,
			"color": fiber.Map{"h": h, "s": s, "v": v},
		})
	}

	// デバッグ用HTTPレス
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok":    true,
		"round": 1,
		"color": fiber.Map{"h": h, "s": s, "v": v},
	})
}

func NextRound(c *fiber.Ctx) error {
	code := c.Params("code")

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	room, err := repositories.GetRoomByCode(ctx, code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"ok": false, "msg": "room not found"})
	}

	newRound, err := repositories.IncrementRound(ctx, room.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": err.Error()})
	}
	if err := repositories.ResetHasCleared(ctx, room.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": err.Error()})
	}

	h, s, v := randomHSV()

	// WS通知：次ラウンド
	if wsMgr != nil {
		hub := wsMgr.GetHub(code)
		hub.Emit("next_round", fiber.Map{
			"round": newRound,
			"color": fiber.Map{"h": h, "s": s, "v": v},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok":    true,
		"round": newRound,
		"color": fiber.Map{"h": h, "s": s, "v": v},
	})
}

func EndRoom(c *fiber.Ctx) error {
	code := c.Params("code")

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	room, err := repositories.GetRoomByCode(ctx, code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"ok": false, "msg": "room not found"})
	}
	if err := repositories.SetRoomEnded(ctx, room.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": err.Error()})
	}

	// WS通知：ゲーム終了
	if wsMgr != nil {
		hub := wsMgr.GetHub(code)
		hub.Emit("game_ended", nil)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// -------- 共通：参加コード生成 & 色乱数 --------

func genCode(n int) string {
	const chars = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"
	rand.Seed(time.Now().UnixNano())
	var b strings.Builder
	b.Grow(n)
	for i := 0; i < n; i++ {
		b.WriteByte(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func randomHSV() (int, float32, float32) {
	rand.Seed(time.Now().UnixNano())
	h := rand.Intn(361)           // 0..360
	s := 0.5 + rand.Float32()*0.4 // 0.5..0.9
	v := 0.7 + rand.Float32()*0.3 // 0.7..1.0
	if v > 1 {
		v = 1
	}
	return h, s, v
}
