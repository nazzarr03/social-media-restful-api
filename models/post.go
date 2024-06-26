package models

import "time"

type Post struct {
	PostID    int       `json:"post_id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	ImageURL  *string   `json:"image_url"`
	Comments  []Comment `json:"comments" gorm:"foreignKey:PostID"`
	Likes     []Like    `json:"likes" gorm:"foreignKey:PostID"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}
