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
	_ = godotenv.Load() // .env があれば読み込む

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// DB接続
	db.Connect()

	// WSマネージャ
	manager := ws.NewManager()

	// ルーター
	app := httpapi.NewRouter(manager)

	// Listen
	log.Printf("🚀 Server running on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
