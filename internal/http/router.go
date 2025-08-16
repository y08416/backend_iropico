package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"backend_iropico/internal/http/handlers"
	"backend_iropico/internal/ws"
)

func NewRouter(hubManager *ws.Manager) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "iropico-backend",
		BodyLimit:    20 * 1024 * 1024, // 20MB: 画像アップロード対策
		ServerHeader: "fiber",
	})

	// ミドルウェア
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // MVPは緩め。必要に応じてフロントのオリジンに絞る
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// ヘルス
	app.Get("/ping", func(c *fiber.Ctx) error { return c.SendString("pong") })

	// Users
	app.Post("/users/register", handlers.RegisterUser)

	// Rooms
	app.Post("/rooms", handlers.CreateRoom)            // ルーム作成（ホスト）
	app.Post("/rooms/:code/join", handlers.JoinRoom)   // 参加
	app.Post("/rooms/:code/start", handlers.StartRoom) // 開始
	app.Post("/rooms/:code/next", handlers.NextRound)  // 次ラウンド
	app.Post("/rooms/:code/end", handlers.EndRoom)     // 終了

	// Submit（写真&スコア）
	app.Post("/rooms/:code/submit", handlers.Submit)

	// Ranking / History
	app.Get("/rooms/:code/ranking", handlers.GetRanking) // ?round=n でラウンド別、無ければ合計
	app.Get("/users/:id/history", handlers.GetUserHistory)

	// WebSocket
	app.Get("/ws/rooms/:code", ws.ServeWS(hubManager)) // ?user_id=&uuid=

	return app
}
