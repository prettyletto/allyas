package storage

import "time"

const SchemaVersion = 1

type Store struct {
	SchemaVersion int     `json:"schema_version"`
	Aliases       []Alias `json:"aliases"`
}

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
