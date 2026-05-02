package service

import (
	"encoding/json"
	"time"

	"github.com/daiyutong/blog/internal/model"
	"gorm.io/gorm"
)

type ProjectService struct {
	db *gorm.DB
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{db: db}
}

func (s *ProjectService) List() ([]model.Project, error) {
	var projects []model.Project
	err := s.db.Order("sort_order ASC, created_at DESC").Find(&projects).Error
	return projects, err
}

func (s *ProjectService) Create(name, description, repoURL, demoURL string, techStack []string, stars int, isFeatured bool) (*model.Project, error) {
	tsJSON, _ := json.Marshal(techStack)
	p := &model.Project{
		Name:        name,
		Description: description,
		RepoURL:     repoURL,
		DemoURL:     demoURL,
		TechStack:   string(tsJSON),
		Stars:       stars,
		IsFeatured:  isFeatured,
		CreatedAt:   time.Now(),
	}
	return p, s.db.Create(p).Error
}

func (s *ProjectService) Update(id uint, name, description, repoURL, demoURL string, techStack []string, stars int, isFeatured bool, sortOrder int) error {
	tsJSON, _ := json.Marshal(techStack)
	return s.db.Model(&model.Project{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name": name, "description": description,
		"repo_url": repoURL, "demo_url": demoURL,
		"tech_stack": string(tsJSON), "stars": stars,
		"is_featured": isFeatured, "sort_order": sortOrder,
	}).Error
}

func (s *ProjectService) Delete(id uint) error {
	return s.db.Delete(&model.Project{}, id).Error
}
