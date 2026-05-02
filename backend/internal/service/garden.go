package service

import (
	"fmt"
	"time"

	"github.com/daiyutong/blog/internal/model"
	"gorm.io/gorm"
)

type GardenService struct {
	db *gorm.DB
}

func NewGardenService(db *gorm.DB) *GardenService {
	return &GardenService{db: db}
}

func (s *GardenService) ListSections() ([]model.GardenSection, error) {
	var sections []model.GardenSection
	err := s.db.Where("is_visible = true AND page = 'garden'").Order("sort_order ASC").Find(&sections).Error
	return sections, err
}

func (s *GardenService) ListWorksSections() ([]model.GardenSection, error) {
	var sections []model.GardenSection
	err := s.db.Where("is_visible = true AND page = 'works'").Order("sort_order ASC").Find(&sections).Error
	return sections, err
}

func (s *GardenService) ListAllSections() ([]model.GardenSection, error) {
	var sections []model.GardenSection
	err := s.db.Order("page ASC, sort_order ASC").Find(&sections).Error
	return sections, err
}

func (s *GardenService) CreateSection(name, slug string) (*model.GardenSection, error) {
	sec := &model.GardenSection{Name: name, Slug: slug, Page: "garden", CreatedAt: time.Now()}
	return sec, s.db.Create(sec).Error
}

func (s *GardenService) CreateWorksSection(name string) (*model.GardenSection, error) {
	// Auto-generate unique slug: works-<timestamp>
	slug := fmt.Sprintf("works-%d", time.Now().UnixNano())
	sec := &model.GardenSection{Name: name, Slug: slug, Page: "works", CreatedAt: time.Now()}
	return sec, s.db.Create(sec).Error
}

func (s *GardenService) UpdateSection(id uint, name string, sortOrder int, isVisible bool) error {
	return s.db.Model(&model.GardenSection{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name": name, "sort_order": sortOrder, "is_visible": isVisible,
	}).Error
}

func (s *GardenService) DeleteSection(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("section_id = ?", id).Delete(&model.GardenItem{})
		return tx.Delete(&model.GardenSection{}, id).Error
	})
}

func (s *GardenService) ListItems(sectionID uint, itemType, tag string) ([]model.GardenItem, error) {
	var items []model.GardenItem
	q := s.db.Where("section_id = ?", sectionID)
	if itemType != "" {
		q = q.Where("type = ?", itemType)
	}
	if tag != "" {
		q = q.Where("tags LIKE ?", "%\""+tag+"\"%")
	}
	err := q.Order("sort_order ASC").Find(&items).Error
	return items, err
}

func (s *GardenService) CreateItem(sectionID uint, itemType, title, url, source, summary, tags, embedURL string) (*model.GardenItem, error) {
	item := &model.GardenItem{
		SectionID: sectionID, Type: itemType, Title: title,
		URL: url, Source: source, Summary: summary, Tags: tags, EmbedURL: embedURL,
		CreatedAt: time.Now(),
	}
	return item, s.db.Create(item).Error
}

func (s *GardenService) UpdateItem(id uint, title, url, source, summary, tags, embedURL string, sortOrder int) error {
	return s.db.Model(&model.GardenItem{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title": title, "url": url, "source": source,
		"summary": summary, "tags": tags, "embed_url": embedURL, "sort_order": sortOrder,
	}).Error
}

func (s *GardenService) DeleteItem(id uint) error {
	return s.db.Delete(&model.GardenItem{}, id).Error
}
