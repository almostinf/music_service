package models

import "time"

type Song struct {
	ID        int64         `json:"id"`
	Title     string        `json:"title"`
	Artist    string        `json:"artist"`
	Duration  time.Duration `json:"duration"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
