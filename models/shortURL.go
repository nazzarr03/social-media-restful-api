package models

import "time"

type ShortURL struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	LongURL   string    `json:"long_url"`
	ShortKey  string    `json:"short_key"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}
