package entities

import "time"

type Image struct {
	ID           uint64       `json:"id"`
	ImageURL     string       `json:"image_url"`
	ThumbnailURL *string      `json:"thumbnail_url,omitempty"`
	Width        *int         `json:"width,omitempty"`
	Height       *int         `json:"height,omitempty"`
	AltText      *string      `json:"alt_text,omitempty"`
	Source       *string      `json:"source,omitempty"`
	IsActive     bool         `json:"is_active"`
	DeletedAt    *time.Time   `json:"deleted_at,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Tags         []TagSummary `json:"tags"`
}

type TagSummary struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CreateImageInput struct {
	ImageURL     string  `json:"image_url"`
	ThumbnailURL *string `json:"thumbnail_url"`
	Width        *int    `json:"width"`
	Height       *int    `json:"height"`
	AltText      *string `json:"alt_text"`
	Source       *string `json:"source"`
}

type UpdateImageInput struct {
	ImageURL     *string `json:"image_url"`
	ThumbnailURL *string `json:"thumbnail_url"`
	Width        *int    `json:"width"`
	Height       *int    `json:"height"`
	AltText      *string `json:"alt_text"`
	Source       *string `json:"source"`
}

type ImageListFilter struct {
	Cursor  *uint64
	Limit   int
	TagSlug string
}
