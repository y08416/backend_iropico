package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend_iropico/internal/db"
	httpapi "backend_iropico/internal/http"
	"backend_iropico/internal/ws"

	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2"
)

func main() {
	_ = godotenv.Load() // .env があれば読み込む

	// 必須ENVの検証（ローカル/本番どちらでも安全）
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL is not set in environment")
	}

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

	// ルート可視化: http://127.0.0.1:3000/_routes で確認できる（開発用）
	app.Get("/_routes", func(c *fiber.Ctx) error {
		return c.JSON(app.Stack())
	})

	// シンプルヘルスチェック（起動確認用）
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	// サーバ起動をゴルーチンで行い、シグナルを待って優雅に停止
	go func() {
		log.Printf("🚀 Server running on :%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatal(err)
		}
	}()

	// Ctrl+C / SIGTERM を待つ
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Println("shutting down...")

	// 優雅に停止
	if err := app.Shutdown(); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
