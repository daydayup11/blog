package service

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/daiyutong/blog/internal/model"
	"gorm.io/gorm"
)

func VisitorHash(ip, userAgent string) string {
	h := sha256.Sum256([]byte(ip + "|" + userAgent))
	return fmt.Sprintf("%x", h[:16])
}

type StatsService struct {
	db *gorm.DB
}

func NewStatsService(db *gorm.DB) *StatsService {
	return &StatsService{db: db}
}

func (s *StatsService) TrackPage(path, visitorHash, referrer string) {
	s.db.Create(&model.PageView{
		Path:        path,
		VisitorHash: visitorHash,
		Referrer:    referrer,
		CreatedAt:   time.Now(),
	})
}

func (s *StatsService) TrackPost(postID uint, visitorHash string) {
	var count int64
	since := time.Now().Add(-24 * time.Hour)
	s.db.Model(&model.PostView{}).
		Where("post_id = ? AND visitor_hash = ? AND created_at > ?", postID, visitorHash, since).
		Count(&count)
	if count == 0 {
		s.db.Create(&model.PostView{
			PostID:      postID,
			VisitorHash: visitorHash,
			CreatedAt:   time.Now(),
		})
	}
}

func (s *StatsService) PostViewCount(postID uint) int64 {
	var count int64
	s.db.Model(&model.PostView{}).Where("post_id = ?", postID).Count(&count)
	return count
}

type DailyStatsResult struct {
	Date string `json:"date"`
	PV   int    `json:"pv"`
	UV   int    `json:"uv"`
}

func (s *StatsService) AdminStats(days int) map[string]interface{} {
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	var rows []DailyStatsResult
	s.db.Model(&model.PageView{}).
		Select("strftime('%Y-%m-%d', created_at) as date, count(*) as pv, count(distinct visitor_hash) as uv").
		Where("strftime('%Y-%m-%d', created_at) >= ?", since).
		Group("date").
		Order("date").
		Scan(&rows)

	var totalPV, totalUV int
	for _, r := range rows {
		totalPV += r.PV
		totalUV += r.UV
	}

	type topPost struct {
		Title string `json:"title"`
		Views int64  `json:"views"`
	}
	var top []topPost
	s.db.Raw(`
		SELECT p.title, COUNT(DISTINCT pv.visitor_hash) as views
		FROM post_views pv
		JOIN posts p ON p.id = pv.post_id
		GROUP BY p.id
		ORDER BY views DESC
		LIMIT 10
	`).Scan(&top)

	return map[string]interface{}{
		"total_pv":  totalPV,
		"total_uv":  totalUV,
		"daily":     rows,
		"top_posts": top,
	}
}
