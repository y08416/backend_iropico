package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend_iropico/internal/db"
	httpapi "backend_iropico/internal/http"
	"backend_iropico/internal/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	// =========================
	// .env の読み込み方針
	// =========================
	// Docker 実行時は docker-compose の environment を信頼し、
	// .env は読み込まない。ローカル実行時のみ .env を読み込む。
	if _, ok := os.LookupEnv("RUNNING_IN_DOCKER"); !ok {
		// 存在するものだけを静かに読む（順に優先度：左→右）
		_ = godotenv.Load("backend/.env", ".env", ".env.localdev")
	}

	// =========================
	// 必須ENVの検証
	// =========================
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL is not set in environment")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// =========================
	// DB接続
	// =========================
	db.Connect()

	// =========================
	// ルーター & WS
	// =========================
	manager := ws.NewManager()
	app := httpapi.NewRouter(manager)

	// 開発用：現在登録済みルートを表示
	app.Get("/_routes", func(c *fiber.Ctx) error {
		return c.JSON(app.Stack())
	})
	// ヘルスチェック
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})
	// ルート
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("iropico-backend is running")
	})

	// =========================
	// 起動 & 優雅な停止
	// =========================
	go func() {
		log.Printf("🚀 Server running on :%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatal(err)
		}
	}()

	// Ctrl+C / SIGTERM を待ってシャットダウン
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Println("shutting down...")

	if err := app.Shutdown(); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
