package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"backend_iropico/internal/repositories"
)

func GetRanking(c *fiber.Ctx) error {
	code := c.Params("code")
	roundStr := c.Query("round")

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	room, err := repositories.GetRoomByCode(ctx, code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"ok": false, "msg": "room not found"})
	}

	if roundStr != "" {
		round, err := strconv.Atoi(roundStr)
		if err != nil || round <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "msg": "invalid round"})
		}
		rows, err := repositories.GetRoundRanking(ctx, room.ID, round)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": err.Error()})
		}
		return c.JSON(fiber.Map{"ok": true, "round": round, "ranking": rows})
	}

	rows, err := repositories.GetTotalRanking(ctx, room.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": err.Error()})
	}
	return c.JSON(fiber.Map{"ok": true, "ranking": rows})
}

func GetUserHistory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || userID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "msg": "invalid user id"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	rows, err := repositories.GetUserHistory(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": err.Error()})
	}
	return c.JSON(fiber.Map{"ok": true, "history": rows})
}
