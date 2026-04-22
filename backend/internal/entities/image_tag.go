package entities

type AttachTagToImageInput struct {
	TagID uint64 `json:"tag_id"`
}

type ImageTagAssignment struct {
	ImageID uint64 `json:"image_id"`
	TagID   uint64 `json:"tag_id"`
}
