package models

import (
	"errors"
	"time"
)

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

type AliasParams struct {
	Group       string
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func NewAlias(id, name, command string, p AliasParams) (*Alias, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if name == "" {
		return nil, errors.New("alias name cannot be empty")
	}
	if command == "" {
		return nil, errors.New("alias command cannot be empty")
	}

	now := time.Now()

	return &Alias{
		ID:          id,
		Name:        name,
		Command:     command,
		Group:       p.Group,
		Description: p.Description,
		Tags:        p.Tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
