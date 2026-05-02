package service

import (
	"time"

	"github.com/daiyutong/blog/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProfileService struct {
	db *gorm.DB
}

func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{db: db}
}

func (s *ProfileService) Get() (*model.Profile, error) {
	var p model.Profile
	err := s.db.First(&p).Error
	if err != nil {
		// Return empty profile if none exists yet
		return &model.Profile{Skills: "[]", Social: "[]", Timeline: "[]"}, nil
	}
	return &p, nil
}

func (s *ProfileService) Save(name, subtitle, bio, skills, social, timeline string) (*model.Profile, error) {
	p := &model.Profile{
		ID:        1,
		Name:      name,
		Subtitle:  subtitle,
		Bio:       bio,
		Skills:    skills,
		Social:    social,
		Timeline:  timeline,
		UpdatedAt: time.Now(),
	}
	err := s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(p).Error
	return p, err
}
