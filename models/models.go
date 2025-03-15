package models

import "time"

type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}
