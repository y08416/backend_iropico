package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5"

	"backend_iropico/internal/db"
	"backend_iropico/internal/models"
)

func main() {
	_ = godotenv.Load(".env")

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
			DisplayName string `json:"display_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.DisplayName == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var u models.User
		err := pool.QueryRow(ctx,
			`INSERT INTO users (display_name)
			 VALUES ($1)
			 RETURNING id, uuid::text, display_name, created_at`,
			in.DisplayName,
		).Scan(&u.ID, &u.UUID, &u.DisplayName, &u.CreatedAt)

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

	log.Println("listening :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}