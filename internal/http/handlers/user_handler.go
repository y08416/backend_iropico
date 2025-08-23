package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"

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

// GET /users/:id
func GetUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	user, err := repositories.GetUserByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, user)
}

// PATCH /users/:id
func UpdateUser(c echo.Context) error {
	// ...実装...
	return nil
}

// DELETE /users/:id
func DeleteUser(c echo.Context) error {
	// ...実装...
	return nil
}

func ListUsers(c echo.Context) error {
	// ユーザー一覧取得の処理をここに実装
	return c.JSON(http.StatusOK, []interface{}{})
}
