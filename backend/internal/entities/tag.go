package entities

import "time"

type Tag struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateTagInput struct {
	Name string  `json:"name"`
	Slug *string `json:"slug"`
}

type UpdateTagInput struct {
	Name *string `json:"name"`
	Slug *string `json:"slug"`
}
