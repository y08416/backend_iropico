package main

import (
	"log"
	"os"

	"backend_iropico/internal/db"
	httpapi "backend_iropico/internal/http"
	"backend_iropico/internal/ws"

	"github.com/joho/godotenv"
)

func main() {
	// .env 読み込み（存在すれば）
	_ = godotenv.Load()

	// Render では PORT 環境変数が与えられる
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// DB 接続
	db.Connect()

	// WebSocket Manager
	manager := ws.NewManager()

	// Fiber Router
	app := httpapi.NewRouter(manager)

	// サーバ起動
	log.Printf("🚀 Server running on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
