package dto

type AddCommentRequest struct {
	Content  string  `json:"content"`
	ImageURL *string `json:"image_url"`
}
