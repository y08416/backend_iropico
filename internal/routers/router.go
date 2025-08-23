package routers

import (
	"backend_iropico/internal/http/handlers"

	"github.com/labstack/echo/v4"
)

func InitRouter(e *echo.Echo) {
	e.GET("/users/:id", handlers.GetUser)
	e.GET("/users", handlers.ListUsers)
	e.PATCH("/users/:id", handlers.UpdateUser)
	e.DELETE("/users/:id", handlers.DeleteUser)
	// 他のエンドポイントも同様に追加
}
