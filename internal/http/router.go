package http

import (
	"backend_iropico/internal/http/handlers"
	"backend_iropico/internal/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)

func NewRouter(hubManager *ws.Manager) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "iropico-backend",
		BodyLimit:    20 * 1024 * 1024, // 20MB: 写真対策
		ServerHeader: "fiber",
	})

	// ★ WSマネージャをハンドラに注入（handlers から Emit できるように）
	handlers.SetWSManager(hubManager)

	// Middlewares
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // ステージングは * でOK。あとで本番ドメインに絞る
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Healthcheck
	app.Get("/ping", func(c *fiber.Ctx) error { return c.SendString("pong") })

	// Users
	app.Post("/users/register", handlers.RegisterUser)
	app.Get("/users/:id/history", handlers.GetUserHistory)

	// Rooms
	app.Post("/rooms", handlers.CreateRoom)            // ルーム作成（ホスト）
	app.Post("/rooms/:code/join", handlers.JoinRoom)   // 参加
	app.Post("/rooms/:code/start", handlers.StartRoom) // 開始
	app.Post("/rooms/:code/next", handlers.NextRound)  // 次ラウンド
	app.Post("/rooms/:code/end", handlers.EndRoom)     // 終了

	// Submit（写真&スコア）
	app.Post("/rooms/:code/submit", handlers.Submit)

	// Ranking
	app.Get("/rooms/:code/ranking", handlers.GetRanking) // ?round=n でラウンド別、無ければ合計

	// WebSocket (upgrade 判定を /ws 配下にだけ適用)
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})
	app.Get("/ws/rooms/:code", websocket.New(ws.ServeWS(hubManager))) // ?user_id=&uuid=

	return app
}
