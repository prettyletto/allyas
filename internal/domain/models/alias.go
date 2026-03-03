package models

import "time"

type Alias struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Command     string    `json:"command"`
	Group       string    `json:"group"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
