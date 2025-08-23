package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"backend_iropico/internal/db"
	"backend_iropico/internal/models"
)

func main() {
	// .env を読むのはローカル実行時のみ（Dockerでは env_file を使用）
	if _, ok := os.LookupEnv("RUNNING_IN_DOCKER"); !ok {
		_ = godotenv.Load("backend/.env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	pool, err := db.NewPoolFromEnv()
	if err != nil {
		log.Fatal("db connect error:", err)
	}
	defer pool.Close()

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var in struct {
			Name string `json:"name"`
			UUID string `json:"uuid"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Name == "" || in.UUID == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var u models.User
		err := pool.QueryRow(ctx,
			`INSERT INTO users (name, uuid)
			 VALUES ($1, $2)
			 ON CONFLICT (uuid) DO UPDATE SET name = EXCLUDED.name
			 RETURNING id, uuid::text, name, created_at`,
			in.Name, in.UUID,
		).Scan(&u.ID, &u.UUID, &u.Name, &u.CreatedAt)

		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "insert failed", http.StatusInternalServerError)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	})

	log.Println("listening :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
