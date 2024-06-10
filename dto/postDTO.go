package dto

type CreatePostRequest struct {
	Content  string  `json:"content"`
	ImageURL *string `json:"image_url"`
}

type EditPostRequest struct {
	Content  string  `json:"content"`
	ImageURL *string `json:"image_url"`
}
