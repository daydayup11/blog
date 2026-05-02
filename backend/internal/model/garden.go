package model

import "time"

type GardenSection struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Slug      string    `gorm:"uniqueIndex;not null" json:"slug"`
	Page      string    `gorm:"default:'garden'" json:"page"` // garden | works
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	IsVisible bool      `gorm:"default:true" json:"is_visible"`
	CreatedAt time.Time `json:"created_at"`
}

type GardenItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SectionID uint      `gorm:"not null" json:"section_id"`
	Type      string    `gorm:"not null" json:"type"`
	Title     string    `gorm:"not null" json:"title"`
	URL       string    `json:"url"`
	Source    string    `json:"source"`
	Summary   string    `json:"summary"`
	Tags      string    `gorm:"default:'[]'" json:"tags"`
	EmbedURL  string    `json:"embed_url"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}
