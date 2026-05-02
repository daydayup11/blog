package model

import "time"

type PostView struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PostID      uint      `gorm:"not null;index" json:"post_id"`
	VisitorHash string    `gorm:"not null" json:"visitor_hash"`
	CreatedAt   time.Time `json:"created_at"`
}

type PageView struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Path        string    `gorm:"not null;index" json:"path"`
	VisitorHash string    `gorm:"not null" json:"visitor_hash"`
	Referrer    string    `json:"referrer"`
	CreatedAt   time.Time `json:"created_at"`
}

type DailyStat struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Date string `gorm:"not null;index" json:"date"`
	Path string `gorm:"not null" json:"path"`
	PV   int    `gorm:"default:0" json:"pv"`
	UV   int    `gorm:"default:0" json:"uv"`
}
