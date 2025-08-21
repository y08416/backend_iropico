package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"backend_iropico/internal/repositories"
)

type registerReq struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

func RegisterUser(c *fiber.Ctx) error {
	var req registerReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "msg": "invalid json"})
	}
	if req.Name == "" || req.UUID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "msg": "name and uuid are required"})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	u, err := repositories.UpsertUserByUUID(ctx, req.Name, req.UUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "msg": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok":   true,
		"user": fiber.Map{"id": u.ID, "name": u.Name, "uuid": u.UUID, "created_at": u.CreatedAt},
	})
}
