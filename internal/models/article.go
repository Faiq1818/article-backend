package models

import "time"

type Article struct {
	ID        string    `json:"id"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
