package models

import "time"

type ShortURL struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	LongURL   string `json:"long_url" gorm:"not null"`
	ShortKey  string `json:"short_key" gorm:"not null"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`

}
