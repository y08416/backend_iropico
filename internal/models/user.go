package models

import "time"

type User struct {
	ID          int64     `json:"id"`
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}
