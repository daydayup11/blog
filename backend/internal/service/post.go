package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/daiyutong/blog/internal/model"
	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

func (s *PostService) Create(title, slug, content, summary string, tags []string) (*model.Post, error) {
	tagsJSON, _ := json.Marshal(tags)
	post := &model.Post{
		Title:     title,
		Slug:      slug,
		Content:   content,
		Summary:   summary,
		WordCount: WordCount(content),
		Tags:      string(tagsJSON),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.db.Create(post).Error; err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}
	return post, nil
}

func (s *PostService) GetBySlug(slug string) (*model.Post, error) {
	var post model.Post
	if err := s.db.Where("slug = ? AND is_published = true", slug).First(&post).Error; err != nil {
		return nil, err
	}
	post.ContentHTML = RenderMarkdown(post.Content)
	return &post, nil
}

// GetAdjacentPosts returns the previous (older) and next (newer) published posts by time.
func (s *PostService) GetAdjacentPosts(createdAt interface{}) (prev, next *model.Post) {
	var p, n model.Post
	// prev = most recent post older than current
	if err := s.db.Where("is_published = true AND created_at < ?", createdAt).
		Order("created_at DESC").First(&p).Error; err == nil {
		prev = &p
	}
	// next = oldest post newer than current
	if err := s.db.Where("is_published = true AND created_at > ?", createdAt).
		Order("created_at ASC").First(&n).Error; err == nil {
		next = &n
	}
	return
}

// GetBySlugAdmin returns a post regardless of published status (for CMS editor).
func (s *PostService) GetBySlugAdmin(slug string) (*model.Post, error) {
	var post model.Post
	if err := s.db.Where("slug = ?", slug).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostService) List(page, pageSize int, tag, query string) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64
	q := s.db.Model(&model.Post{}).Where("is_published = true")
	if tag != "" {
		q = q.Where("tags LIKE ?", "%\""+tag+"\"%")
	}
	if query != "" {
		q = q.Where("title LIKE ? OR summary LIKE ?", "%"+query+"%", "%"+query+"%")
	}
	q.Count(&total)
	offset := (page - 1) * pageSize
	err := q.Order("is_pinned DESC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&posts).Error
	return posts, total, err
}

// ListAll returns all posts regardless of published status (admin only).
func (s *PostService) ListAll(page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64
	s.db.Model(&model.Post{}).Count(&total)
	offset := (page - 1) * pageSize
	err := s.db.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&posts).Error
	return posts, total, err
}

func (s *PostService) Update(id uint, title, content, summary string, tags []string) (*model.Post, error) {
	tagsJSON, _ := json.Marshal(tags)
	updates := map[string]interface{}{
		"title":      title,
		"content":    content,
		"summary":    summary,
		"word_count": WordCount(content),
		"tags":       string(tagsJSON),
		"updated_at": time.Now(),
	}
	if err := s.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	var post model.Post
	s.db.First(&post, id)
	return &post, nil
}

func (s *PostService) Delete(id uint) error {
	return s.db.Delete(&model.Post{}, id).Error
}

func (s *PostService) SetPinned(id uint, pinned bool) error {
	return s.db.Model(&model.Post{}).Where("id = ?", id).Update("is_pinned", pinned).Error
}

func (s *PostService) SetPublished(id uint, published bool) error {
	return s.db.Model(&model.Post{}).Where("id = ?", id).Update("is_published", published).Error
}

func (s *PostService) AllTags() []string {
	var posts []model.Post
	s.db.Select("tags").Where("is_published = true").Find(&posts)
	seen := map[string]bool{}
	var tags []string
	for _, p := range posts {
		var t []string
		json.Unmarshal([]byte(p.Tags), &t)
		for _, tag := range t {
			if !seen[tag] {
				seen[tag] = true
				tags = append(tags, tag)
			}
		}
	}
	if tags == nil {
		return []string{}
	}
	return tags
}
