package models

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	UUID      string `json:"uuid"`
	CreatedAt string `json:"created_at"`
}
