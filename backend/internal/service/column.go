package service

import (
	"fmt"
	"time"

	"github.com/daiyutong/blog/internal/model"
	"gorm.io/gorm"
)

type ColumnService struct {
	db *gorm.DB
}

func NewColumnService(db *gorm.DB) *ColumnService {
	return &ColumnService{db: db}
}

func (s *ColumnService) Create(name, slug, description, coverURL string) (*model.Column, error) {
	col := &model.Column{
		Name:        name,
		Slug:        slug,
		Description: description,
		CoverURL:    coverURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.db.Create(col).Error; err != nil {
		return nil, fmt.Errorf("create column: %w", err)
	}
	return col, nil
}

func (s *ColumnService) List() ([]model.Column, error) {
	var cols []model.Column
	err := s.db.Where("is_published = true").Order("created_at DESC").Find(&cols).Error
	return cols, err
}

func (s *ColumnService) ListAll() ([]model.Column, error) {
	var cols []model.Column
	err := s.db.Order("created_at DESC").Find(&cols).Error
	return cols, err
}

func (s *ColumnService) GetBySlug(slug string) (*model.Column, error) {
	var col model.Column
	if err := s.db.Where("slug = ? AND is_published = true", slug).First(&col).Error; err != nil {
		return nil, err
	}
	return &col, nil
}

func (s *ColumnService) SetPosts(columnID uint, postIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("column_id = ?", columnID).Delete(&model.ColumnPost{})
		for i, pid := range postIDs {
			cp := model.ColumnPost{ColumnID: columnID, PostID: pid, SortOrder: i}
			if err := tx.Create(&cp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *ColumnService) GetPosts(columnID uint) ([]model.Post, error) {
	var posts []model.Post
	err := s.db.
		Joins("JOIN column_posts cp ON cp.post_id = posts.id").
		Where("cp.column_id = ?", columnID).
		Order("cp.sort_order ASC").
		Find(&posts).Error
	return posts, err
}

func (s *ColumnService) Update(id uint, name, description, coverURL string) error {
	return s.db.Model(&model.Column{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":        name,
		"description": description,
		"cover_url":   coverURL,
		"updated_at":  time.Now(),
	}).Error
}

func (s *ColumnService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("column_id = ?", id).Delete(&model.ColumnPost{})
		return tx.Delete(&model.Column{}, id).Error
	})
}

func (s *ColumnService) SetPublished(id uint, published bool) error {
	return s.db.Model(&model.Column{}).Where("id = ?", id).Update("is_published", published).Error
}
