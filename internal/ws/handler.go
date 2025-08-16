package ws

import "github.com/gofiber/fiber/v2"

func ServeWS(m *Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"ok":  false,
			"msg": "WebSocket endpoint not implemented yet",
		})
	}
}
