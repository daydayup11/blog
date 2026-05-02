package model

import "time"

// Profile stores the site owner's personal information shown on the homepage.
// Only one row exists (id=1); upsert on save.
type Profile struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Subtitle  string    `json:"subtitle"`
	Bio       string    `json:"bio"`
	Skills    string    `gorm:"default:'[]'" json:"skills"`    // JSON string array
	Social    string    `gorm:"default:'[]'" json:"social"`    // JSON [{label,url,icon}]
	Timeline  string    `gorm:"default:'[]'" json:"timeline"` // JSON [{year,title,sub}]
	UpdatedAt time.Time `json:"updated_at"`
}
