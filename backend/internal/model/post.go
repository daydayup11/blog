package model

import "time"

type Post struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Slug        string    `gorm:"uniqueIndex;not null" json:"slug"`
	Content     string    `gorm:"not null" json:"-"`
	ContentHTML string    `gorm:"-" json:"content_html"`
	Summary     string    `json:"summary"`
	WordCount   int       `json:"word_count"`
	Tags        string    `gorm:"default:'[]'" json:"tags"`
	IsPinned    bool      `gorm:"default:false" json:"is_pinned"`
	IsPublished bool      `gorm:"default:false" json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Column struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Slug        string    `gorm:"uniqueIndex;not null" json:"slug"`
	CoverURL    string    `json:"cover_url"`
	Description string    `json:"description"`
	IsPublished bool      `gorm:"default:false" json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ColumnPost struct {
	ColumnID  uint `gorm:"primaryKey" json:"column_id"`
	PostID    uint `gorm:"primaryKey" json:"post_id"`
	SortOrder int  `gorm:"not null" json:"sort_order"`
}
