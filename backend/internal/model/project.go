package model

import "time"

type Project struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	RepoURL     string    `json:"repo_url"`
	DemoURL     string    `json:"demo_url"`
	TechStack   string    `gorm:"default:'[]'" json:"tech_stack"`
	Stars       int       `gorm:"default:0" json:"stars"`
	IsFeatured  bool      `gorm:"default:false" json:"is_featured"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}
