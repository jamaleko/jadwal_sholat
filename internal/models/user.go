package models

import "time"

type User struct {
	ChatID    int64     `db:"chat_id"`
	Latitude  float64   `db:"latitude"`
	Longitude float64   `db:"longitude"`
	Active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
}
