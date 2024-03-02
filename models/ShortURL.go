package models

import ()

type ShortURL struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	LongURL   string `json:"long_url" gorm:"not null"`
	ShortKey  string `json:"short_key" gorm:"not null"`
	CreatedAt string `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt string `json:"updated_at" gorm:"default:current_timestamp"`
}
