package models

import "time"

type Event struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Date      time.Time `json:"event_date"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
