# Backend API Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the complete Go REST API backend: DB schema, all public and admin endpoints, JWT auth, analytics tracking, and Markdown processing.

**Architecture:** Gin HTTP framework with GORM + SQLite. Code split by responsibility: `handler/` for HTTP logic, `service/` for business logic, `model/` for DB models. JWT middleware guards all `/api/admin/*` routes. A daily cron goroutine aggregates `page_views` into `daily_stats` at midnight.

**Tech Stack:** Go 1.22+, Gin, GORM, SQLite (gorm driver), goldmark (Markdown→HTML + word count), golang-jwt/jwt v5, crypto/sha256 (visitor hash)

---

## File Structure

```
backend/
├── cmd/server/
│   └── main.go                  # Entry point: wire everything, start server
├── internal/
│   ├── config/
│   │   └── config.go            # Env var loading (DB_PATH, JWT_SECRET, PORT)
│   ├── db/
│   │   └── db.go                # GORM init, AutoMigrate all models
│   ├── model/
│   │   ├── post.go              # Post, Tag, Column, ColumnPost models
│   │   ├── project.go           # Project model
│   │   ├── garden.go            # GardenSection, GardenItem models
│   │   └── stats.go             # PostView, PageView, DailyStat models
│   ├── service/
│   │   ├── post.go              # Post CRUD, word count, pinned logic
│   │   ├── column.go            # Column CRUD, ordered post list
│   │   ├── project.go           # Project CRUD
│   │   ├── garden.go            # Section + Item CRUD
│   │   ├── stats.go             # Track page/post view, aggregate daily stats
│   │   └── upload.go            # Save uploaded image to ./data/uploads/
│   ├── handler/
│   │   ├── post.go              # GET /api/posts, GET /api/posts/:slug, admin CRUD
│   │   ├── column.go            # GET /api/columns, GET /api/columns/:slug, admin CRUD
│   │   ├── project.go           # GET /api/projects, admin CRUD
│   │   ├── garden.go            # GET /api/garden/*, admin CRUD
│   │   ├── stats.go             # POST /api/track/*, GET /api/admin/stats
│   │   ├── auth.go              # POST /api/auth/login
│   │   └── upload.go            # POST /api/admin/upload
│   ├── middleware/
│   │   ├── auth.go              # JWT validation middleware
│   │   └── cors.go              # CORS headers for local dev
│   └── router/
│       └── router.go            # All route registrations
├── go.mod
├── go.sum
└── Dockerfile
```

**Test files** (mirror source structure):
```
backend/
├── internal/service/post_test.go
├── internal/service/column_test.go
├── internal/service/stats_test.go
└── internal/handler/post_test.go   # HTTP integration tests via httptest
```

---

## Task 1: Go Module + Dependencies

**Files:**
- Create: `backend/go.mod`
- Create: `backend/go.sum` (auto-generated)

- [ ] **Step 1: Initialize Go module**

```bash
cd /Users/daiyutong/IdeaProjects/blog/backend
go mod init github.com/daiyutong/blog
```

- [ ] **Step 2: Install dependencies**

```bash
go get github.com/gin-gonic/gin@v1.10.0
go get gorm.io/gorm@v1.25.10
go get gorm.io/driver/sqlite@v1.5.5
go get github.com/golang-jwt/jwt/v5@v5.2.1
go get github.com/yuin/goldmark@v1.7.1
go get github.com/yuin/goldmark-highlighting/v2@v2.0.0-20230729083705-37449abec8cc
go get github.com/alecthomas/chroma/v2@v2.13.0
go tidy
```

- [ ] **Step 3: Commit**

```bash
git add backend/go.mod backend/go.sum
git commit -m "feat: init Go module and dependencies"
```

---

## Task 2: Config + DB Init

**Files:**
- Create: `backend/internal/config/config.go`
- Create: `backend/internal/db/db.go`

- [ ] **Step 1: Write config loader**

Create `backend/internal/config/config.go`:

```go
package config

import "os"

type Config struct {
	DBPath    string
	JWTSecret string
	Port      string
	AdminUser string
	AdminPass string
}

func Load() Config {
	return Config{
		DBPath:    getEnv("DB_PATH", "./data/blog.db"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-in-prod"),
		Port:      getEnv("PORT", "8080"),
		AdminUser: getEnv("ADMIN_USER", "admin"),
		AdminPass: getEnv("ADMIN_PASS", "admin123"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
```

- [ ] **Step 2: Write DB init**

Create `backend/internal/db/db.go`:

```go
package db

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(dbPath string) *gorm.DB {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("failed to create db dir: %v", err)
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	return db
}
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/
git commit -m "feat: add config loader and db init"
```

---

## Task 3: Data Models

**Files:**
- Create: `backend/internal/model/post.go`
- Create: `backend/internal/model/project.go`
- Create: `backend/internal/model/garden.go`
- Create: `backend/internal/model/stats.go`

- [ ] **Step 1: Write post models**

Create `backend/internal/model/post.go`:

```go
package model

import "time"

type Post struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Slug        string    `gorm:"uniqueIndex;not null" json:"slug"`
	Content     string    `gorm:"not null" json:"-"`         // raw Markdown, omit from list API
	ContentHTML string    `gorm:"-" json:"content_html"`     // rendered, not stored
	Summary     string    `json:"summary"`
	WordCount   int       `json:"word_count"`
	Tags        string    `gorm:"default:'[]'" json:"tags"`  // JSON array string
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
```

- [ ] **Step 2: Write project model**

Create `backend/internal/model/project.go`:

```go
package model

import "time"

type Project struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	RepoURL     string    `json:"repo_url"`
	DemoURL     string    `json:"demo_url"`
	TechStack   string    `gorm:"default:'[]'" json:"tech_stack"` // JSON array string
	Stars       int       `gorm:"default:0" json:"stars"`
	IsFeatured  bool      `gorm:"default:false" json:"is_featured"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}
```

- [ ] **Step 3: Write garden models**

Create `backend/internal/model/garden.go`:

```go
package model

import "time"

type GardenSection struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Slug      string    `gorm:"uniqueIndex;not null" json:"slug"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	IsVisible bool      `gorm:"default:true" json:"is_visible"`
	CreatedAt time.Time `json:"created_at"`
}

type GardenItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SectionID uint      `gorm:"not null" json:"section_id"`
	Type      string    `gorm:"not null" json:"type"` // article/video/book/account/other
	Title     string    `gorm:"not null" json:"title"`
	URL       string    `json:"url"`
	Source    string    `json:"source"`
	Summary   string    `json:"summary"`
	EmbedURL  string    `json:"embed_url"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}
```

- [ ] **Step 4: Write stats models**

Create `backend/internal/model/stats.go`:

```go
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
	Date string `gorm:"not null;index" json:"date"` // YYYY-MM-DD
	Path string `gorm:"not null" json:"path"`
	PV   int    `gorm:"default:0" json:"pv"`
	UV   int    `gorm:"default:0" json:"uv"`
}
```

- [ ] **Step 5: Wire AutoMigrate into db.go**

Update `backend/internal/db/db.go` — add import and call after `gorm.Open`:

```go
package db

import (
	"log"
	"os"
	"path/filepath"

	"github.com/daiyutong/blog/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(dbPath string) *gorm.DB {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("failed to create db dir: %v", err)
	}
	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	if err := database.AutoMigrate(
		&model.Post{},
		&model.Column{},
		&model.ColumnPost{},
		&model.Project{},
		&model.GardenSection{},
		&model.GardenItem{},
		&model.PostView{},
		&model.PageView{},
		&model.DailyStat{},
	); err != nil {
		log.Fatalf("automigrate failed: %v", err)
	}
	return database
}
```

- [ ] **Step 6: Commit**

```bash
git add backend/internal/model/ backend/internal/db/
git commit -m "feat: add all GORM models and automigrate"
```

---

## Task 4: Markdown Service (word count + HTML render)

**Files:**
- Create: `backend/internal/service/markdown.go`
- Create: `backend/internal/service/markdown_test.go`

- [ ] **Step 1: Write failing tests**

Create `backend/internal/service/markdown_test.go`:

```go
package service_test

import (
	"testing"

	"github.com/daiyutong/blog/internal/service"
)

func TestWordCount_ChineseAndEnglish(t *testing.T) {
	md := "Hello world 你好世界"
	count := service.WordCount(md)
	// 2 English words + 4 Chinese chars = 6
	if count != 6 {
		t.Errorf("expected 6, got %d", count)
	}
}

func TestWordCount_StripsMarkdown(t *testing.T) {
	md := "# 标题\n\n**加粗** 文字"
	count := service.WordCount(md)
	// "标题" + "加粗" + "文字" = 6 Chinese chars
	if count != 6 {
		t.Errorf("expected 6, got %d", count)
	}
}

func TestRenderMarkdown_ReturnsHTML(t *testing.T) {
	md := "# Hello\n\nworld"
	html := service.RenderMarkdown(md)
	if html == "" {
		t.Error("expected non-empty HTML")
	}
	if len(html) < 10 {
		t.Errorf("HTML too short: %s", html)
	}
}

func TestReadingMinutes(t *testing.T) {
	if service.ReadingMinutes(0) != 1 {
		t.Error("0 words should return 1 min minimum")
	}
	if service.ReadingMinutes(300) != 1 {
		t.Error("300 words should return 1 min")
	}
	if service.ReadingMinutes(600) != 2 {
		t.Error("600 words should return 2 min")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/daiyutong/IdeaProjects/blog/backend
go test ./internal/service/... -v -run "TestWordCount|TestRenderMarkdown|TestReadingMinutes"
```

Expected: `cannot find package` or `undefined: service.WordCount`

- [ ] **Step 3: Implement markdown service**

Create `backend/internal/service/markdown.go`:

```go
package service

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM, extension.Table),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	goldmark.WithRendererOptions(html.WithHardWraps(), html.WithUnsafe()),
)

var markdownSymbols = regexp.MustCompile(`[#*_\[\]()~` + "`" + `>|!\-=]+`)

func RenderMarkdown(source string) string {
	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		return source
	}
	return buf.String()
}

func WordCount(source string) int {
	// Strip markdown syntax
	plain := markdownSymbols.ReplaceAllString(source, " ")
	count := 0
	inWord := false
	for _, r := range plain {
		if unicode.Is(unicode.Han, r) {
			count++
			inWord = false
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if !inWord {
				count++
				inWord = true
			}
		} else {
			inWord = false
		}
	}
	return count
}

func ReadingMinutes(wordCount int) int {
	mins := wordCount / 300
	if mins < 1 {
		return 1
	}
	return mins
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd /Users/daiyutong/IdeaProjects/blog/backend
go test ./internal/service/... -v -run "TestWordCount|TestRenderMarkdown|TestReadingMinutes"
```

Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add backend/internal/service/markdown.go backend/internal/service/markdown_test.go
git commit -m "feat: markdown render and word count service"
```

---

## Task 5: Stats Service (visitor hash + tracking)

**Files:**
- Create: `backend/internal/service/stats.go`
- Create: `backend/internal/service/stats_test.go`

- [ ] **Step 1: Write failing tests**

Create `backend/internal/service/stats_test.go`:

```go
package service_test

import (
	"testing"

	"github.com/daiyutong/blog/internal/service"
)

func TestVisitorHash_Deterministic(t *testing.T) {
	h1 := service.VisitorHash("192.168.1.1", "Mozilla/5.0")
	h2 := service.VisitorHash("192.168.1.1", "Mozilla/5.0")
	if h1 != h2 {
		t.Error("same inputs should produce same hash")
	}
}

func TestVisitorHash_Different(t *testing.T) {
	h1 := service.VisitorHash("192.168.1.1", "Mozilla/5.0")
	h2 := service.VisitorHash("10.0.0.1", "Mozilla/5.0")
	if h1 == h2 {
		t.Error("different IPs should produce different hashes")
	}
}

func TestVisitorHash_NotRawIP(t *testing.T) {
	h := service.VisitorHash("192.168.1.1", "agent")
	if h == "192.168.1.1" {
		t.Error("hash must not be the raw IP")
	}
	if len(h) < 16 {
		t.Error("hash too short")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/daiyutong/IdeaProjects/blog/backend
go test ./internal/service/... -v -run "TestVisitorHash"
```

Expected: FAIL with `undefined: service.VisitorHash`

- [ ] **Step 3: Implement stats service**

Create `backend/internal/service/stats.go`:

```go
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
	// Only count once per visitor per post per 24h
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
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd /Users/daiyutong/IdeaProjects/blog/backend
go test ./internal/service/... -v -run "TestVisitorHash"
```

Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add backend/internal/service/stats.go backend/internal/service/stats_test.go
git commit -m "feat: stats service with visitor hash and tracking"
```

---

## Task 6: Post Service

**Files:**
- Create: `backend/internal/service/post.go`
- Create: `backend/internal/service/post_test.go`

- [ ] **Step 1: Write failing tests**

Create `backend/internal/service/post_test.go`:

```go
package service_test

import (
	"testing"

	"github.com/daiyutong/blog/internal/db"
	"github.com/daiyutong/blog/internal/service"
)

func setupTestDB(t *testing.T) *service.PostService {
	t.Helper()
	database := db.Init(":memory:")
	return service.NewPostService(database)
}

func TestCreatePost_SetsWordCount(t *testing.T) {
	svc := setupTestDB(t)
	post, err := svc.Create("测试标题", "test-slug", "你好世界 hello", "摘要", []string{"go"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	// "你好世界" = 4 chars, "hello" = 1 word = 5
	if post.WordCount != 5 {
		t.Errorf("expected word count 5, got %d", post.WordCount)
	}
}

func TestCreatePost_SlugUnique(t *testing.T) {
	svc := setupTestDB(t)
	_, err := svc.Create("标题1", "same-slug", "内容", "", []string{})
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}
	_, err = svc.Create("标题2", "same-slug", "内容2", "", []string{})
	if err == nil {
		t.Error("expected error for duplicate slug")
	}
}

func TestSetPinned(t *testing.T) {
	svc := setupTestDB(t)
	post, _ := svc.Create("标题", "pin-slug", "内容", "", []string{})
	if err := svc.SetPinned(post.ID, true); err != nil {
		t.Fatalf("set pinned failed: %v", err)
	}
	fetched, _ := svc.GetBySlug("pin-slug")
	if !fetched.IsPinned {
		t.Error("expected post to be pinned")
	}
}

func TestListPosts_PinnedFirst(t *testing.T) {
	svc := setupTestDB(t)
	svc.Create("普通", "normal", "内容", "", []string{})
	p2, _ := svc.Create("置顶", "pinned", "内容", "", []string{})
	svc.SetPinned(p2.ID, true)
	svc.SetPublished(p2.ID, true)
	svc.SetPublished(1, true)

	posts, _, _ := svc.List(1, 10, "", "")
	if len(posts) == 0 {
		t.Fatal("expected posts")
	}
	if !posts[0].IsPinned {
		t.Error("first post should be pinned")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/daiyutong/IdeaProjects/blog/backend
go test ./internal/service/... -v -run "TestCreatePost|TestSetPinned|TestListPosts"
```

Expected: FAIL

- [ ] **Step 3: Implement post service**

Create `backend/internal/service/post.go`:

```go
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

func (s *PostService) List(page, pageSize int, tag, query string) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64
	db := s.db.Model(&model.Post{}).Where("is_published = true")
	if tag != "" {
		db = db.Where("tags LIKE ?", "%\""+tag+"\"%")
	}
	if query != "" {
		db = db.Where("title LIKE ? OR summary LIKE ?", "%"+query+"%", "%"+query+"%")
	}
	db.Count(&total)
	offset := (page - 1) * pageSize
	err := db.Order("is_pinned DESC, created_at DESC").
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
	return tags
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd /Users/daiyutong/IdeaProjects/blog/backend
go test ./internal/service/... -v -run "TestCreatePost|TestSetPinned|TestListPosts"
```

Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add backend/internal/service/post.go backend/internal/service/post_test.go
git commit -m "feat: post service with CRUD, pinning, word count"
```

---

## Task 7: Column Service

**Files:**
- Create: `backend/internal/service/column.go`
- Create: `backend/internal/service/column_test.go`

- [ ] **Step 1: Write failing tests**

Create `backend/internal/service/column_test.go`:

```go
package service_test

import (
	"testing"

	"github.com/daiyutong/blog/internal/db"
	"github.com/daiyutong/blog/internal/service"
)

func setupColumnDB(t *testing.T) (*service.ColumnService, *service.PostService) {
	t.Helper()
	database := db.Init(":memory:")
	return service.NewColumnService(database), service.NewPostService(database)
}

func TestCreateColumn(t *testing.T) {
	svc, _ := setupColumnDB(t)
	col, err := svc.Create("Go 系列", "go-series", "Go 从入门到实践", "")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if col.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestColumnPostOrder(t *testing.T) {
	colSvc, postSvc := setupColumnDB(t)
	col, _ := colSvc.Create("系列", "series", "", "")
	p1, _ := postSvc.Create("第一篇", "post-1", "内容", "", []string{})
	p2, _ := postSvc.Create("第二篇", "post-2", "内容", "", []string{})

	err := colSvc.SetPosts(col.ID, []uint{p1.ID, p2.ID})
	if err != nil {
		t.Fatalf("set posts failed: %v", err)
	}

	posts, err := colSvc.GetPosts(col.ID)
	if err != nil {
		t.Fatalf("get posts failed: %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
	if posts[0].ID != p1.ID {
		t.Error("first post should be p1")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/daiyutong/IdeaProjects/blog/backend
go test ./internal/service/... -v -run "TestCreateColumn|TestColumnPostOrder"
```

Expected: FAIL

- [ ] **Step 3: Implement column service**

Create `backend/internal/service/column.go`:

```go
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

func (s *ColumnService) GetBySlug(slug string) (*model.Column, error) {
	var col model.Column
	if err := s.db.Where("slug = ? AND is_published = true", slug).First(&col).Error; err != nil {
		return nil, err
	}
	return &col, nil
}

// SetPosts replaces the ordered post list for a column.
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
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd /Users/daiyutong/IdeaProjects/blog/backend
go test ./internal/service/... -v -run "TestCreateColumn|TestColumnPostOrder"
```

Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add backend/internal/service/column.go backend/internal/service/column_test.go
git commit -m "feat: column service with ordered post management"
```

---

## Task 8: Project + Garden Services

**Files:**
- Create: `backend/internal/service/project.go`
- Create: `backend/internal/service/garden.go`

- [ ] **Step 1: Implement project service**

Create `backend/internal/service/project.go`:

```go
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
```

- [ ] **Step 2: Implement garden service**

Create `backend/internal/service/garden.go`:

```go
package service

import (
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
	err := s.db.Where("is_visible = true").Order("sort_order ASC").Find(&sections).Error
	return sections, err
}

func (s *GardenService) CreateSection(name, slug string) (*model.GardenSection, error) {
	sec := &model.GardenSection{Name: name, Slug: slug, CreatedAt: time.Now()}
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

func (s *GardenService) ListItems(sectionID uint, itemType string) ([]model.GardenItem, error) {
	var items []model.GardenItem
	q := s.db.Where("section_id = ?", sectionID)
	if itemType != "" {
		q = q.Where("type = ?", itemType)
	}
	err := q.Order("sort_order ASC").Find(&items).Error
	return items, err
}

func (s *GardenService) CreateItem(sectionID uint, itemType, title, url, source, summary, embedURL string) (*model.GardenItem, error) {
	item := &model.GardenItem{
		SectionID: sectionID, Type: itemType, Title: title,
		URL: url, Source: source, Summary: summary, EmbedURL: embedURL,
		CreatedAt: time.Now(),
	}
	return item, s.db.Create(item).Error
}

func (s *GardenService) UpdateItem(id uint, title, url, source, summary, embedURL string, sortOrder int) error {
	return s.db.Model(&model.GardenItem{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title": title, "url": url, "source": source,
		"summary": summary, "embed_url": embedURL, "sort_order": sortOrder,
	}).Error
}

func (s *GardenService) DeleteItem(id uint) error {
	return s.db.Delete(&model.GardenItem{}, id).Error
}
```

- [ ] **Step 3: Run all service tests**

```bash
cd /Users/daiyutong/IdeaProjects/blog/backend
go test ./internal/service/... -v
```

Expected: all PASS

- [ ] **Step 4: Commit**

```bash
git add backend/internal/service/project.go backend/internal/service/garden.go
git commit -m "feat: project and garden services"
```

---

## Task 9: JWT Auth Middleware

**Files:**
- Create: `backend/internal/middleware/auth.go`
- Create: `backend/internal/middleware/cors.go`

- [ ] **Step 1: Implement JWT middleware**

Create `backend/internal/middleware/auth.go`:

```go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Next()
	}
}
```

- [ ] **Step 2: Implement CORS middleware**

Create `backend/internal/middleware/cors.go`:

```go
package middleware

import (
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/middleware/
git commit -m "feat: JWT auth and CORS middleware"
```

---

## Task 10: HTTP Handlers

**Files:**
- Create: `backend/internal/handler/auth.go`
- Create: `backend/internal/handler/post.go`
- Create: `backend/internal/handler/column.go`
- Create: `backend/internal/handler/project.go`
- Create: `backend/internal/handler/garden.go`
- Create: `backend/internal/handler/stats.go`
- Create: `backend/internal/handler/upload.go`

- [ ] **Step 1: Auth handler**

Create `backend/internal/handler/auth.go`:

```go
package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	adminUser string
	adminPass string
	jwtSecret string
}

func NewAuthHandler(adminUser, adminPass, jwtSecret string) *AuthHandler {
	return &AuthHandler{adminUser: adminUser, adminPass: adminPass, jwtSecret: jwtSecret}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.Username != h.adminUser || req.Password != h.adminPass {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": req.Username,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token signing failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": signed})
}
```

- [ ] **Step 2: Post handler**

Create `backend/internal/handler/post.go`:

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/daiyutong/blog/internal/service"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	posts *service.PostService
	stats *service.StatsService
}

func NewPostHandler(posts *service.PostService, stats *service.StatsService) *PostHandler {
	return &PostHandler{posts: posts, stats: stats}
}

func (h *PostHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	tag := c.Query("tag")
	q := c.Query("q")
	posts, total, err := h.posts.List(page, 10, tag, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts, "total": total, "page": page})
}

func (h *PostHandler) GetBySlug(c *gin.Context) {
	post, err := h.posts.GetBySlug(c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	views := h.stats.PostViewCount(post.ID)
	c.JSON(http.StatusOK, gin.H{"post": post, "views": views, "reading_minutes": service.ReadingMinutes(post.WordCount)})
}

func (h *PostHandler) Tags(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tags": h.posts.AllTags()})
}

func (h *PostHandler) AdminCreate(c *gin.Context) {
	var req struct {
		Title   string   `json:"title" binding:"required"`
		Slug    string   `json:"slug" binding:"required"`
		Content string   `json:"content" binding:"required"`
		Summary string   `json:"summary"`
		Tags    []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post, err := h.posts.Create(req.Title, req.Slug, req.Content, req.Summary, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) AdminUpdate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Summary string   `json:"summary"`
		Tags    []string `json:"tags"`
	}
	c.ShouldBindJSON(&req)
	post, err := h.posts.Update(uint(id), req.Title, req.Content, req.Summary, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) AdminDelete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.posts.Delete(uint(id))
	c.Status(http.StatusNoContent)
}

func (h *PostHandler) AdminSetPin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct{ Pinned bool `json:"pinned"` }
	c.ShouldBindJSON(&req)
	h.posts.SetPinned(uint(id), req.Pinned)
	c.Status(http.StatusNoContent)
}

func (h *PostHandler) AdminSetPublished(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct{ Published bool `json:"published"` }
	c.ShouldBindJSON(&req)
	h.posts.SetPublished(uint(id), req.Published)
	c.Status(http.StatusNoContent)
}
```

- [ ] **Step 3: Column handler**

Create `backend/internal/handler/column.go`:

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/daiyutong/blog/internal/service"
	"github.com/gin-gonic/gin"
)

type ColumnHandler struct {
	columns *service.ColumnService
}

func NewColumnHandler(columns *service.ColumnService) *ColumnHandler {
	return &ColumnHandler{columns: columns}
}

func (h *ColumnHandler) List(c *gin.Context) {
	cols, err := h.columns.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"columns": cols})
}

func (h *ColumnHandler) GetBySlug(c *gin.Context) {
	col, err := h.columns.GetBySlug(c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	posts, _ := h.columns.GetPosts(col.ID)
	c.JSON(http.StatusOK, gin.H{"column": col, "posts": posts})
}

func (h *ColumnHandler) AdminCreate(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Slug        string `json:"slug" binding:"required"`
		Description string `json:"description"`
		CoverURL    string `json:"cover_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	col, err := h.columns.Create(req.Name, req.Slug, req.Description, req.CoverURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, col)
}

func (h *ColumnHandler) AdminUpdate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		CoverURL    string `json:"cover_url"`
	}
	c.ShouldBindJSON(&req)
	h.columns.Update(uint(id), req.Name, req.Description, req.CoverURL)
	c.Status(http.StatusNoContent)
}

func (h *ColumnHandler) AdminDelete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.columns.Delete(uint(id))
	c.Status(http.StatusNoContent)
}

func (h *ColumnHandler) AdminSetPosts(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct{ PostIDs []uint `json:"post_ids"` }
	c.ShouldBindJSON(&req)
	if err := h.columns.SetPosts(uint(id), req.PostIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
```

- [ ] **Step 4: Project, Garden, Stats, Upload handlers**

Create `backend/internal/handler/project.go`:

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/daiyutong/blog/internal/service"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct{ projects *service.ProjectService }

func NewProjectHandler(p *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projects: p}
}

func (h *ProjectHandler) List(c *gin.Context) {
	projects, err := h.projects.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (h *ProjectHandler) AdminCreate(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		RepoURL     string   `json:"repo_url"`
		DemoURL     string   `json:"demo_url"`
		TechStack   []string `json:"tech_stack"`
		Stars       int      `json:"stars"`
		IsFeatured  bool     `json:"is_featured"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := h.projects.Create(req.Name, req.Description, req.RepoURL, req.DemoURL, req.TechStack, req.Stars, req.IsFeatured)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *ProjectHandler) AdminUpdate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		RepoURL     string   `json:"repo_url"`
		DemoURL     string   `json:"demo_url"`
		TechStack   []string `json:"tech_stack"`
		Stars       int      `json:"stars"`
		IsFeatured  bool     `json:"is_featured"`
		SortOrder   int      `json:"sort_order"`
	}
	c.ShouldBindJSON(&req)
	h.projects.Update(uint(id), req.Name, req.Description, req.RepoURL, req.DemoURL, req.TechStack, req.Stars, req.IsFeatured, req.SortOrder)
	c.Status(http.StatusNoContent)
}

func (h *ProjectHandler) AdminDelete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.projects.Delete(uint(id))
	c.Status(http.StatusNoContent)
}
```

Create `backend/internal/handler/garden.go`:

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/daiyutong/blog/internal/service"
	"github.com/gin-gonic/gin"
)

type GardenHandler struct{ garden *service.GardenService }

func NewGardenHandler(g *service.GardenService) *GardenHandler {
	return &GardenHandler{garden: g}
}

func (h *GardenHandler) ListSections(c *gin.Context) {
	sections, _ := h.garden.ListSections()
	c.JSON(http.StatusOK, gin.H{"sections": sections})
}

func (h *GardenHandler) ListItems(c *gin.Context) {
	sectionID, _ := strconv.ParseUint(c.Query("section_id"), 10, 64)
	items, _ := h.garden.ListItems(uint(sectionID), c.Query("type"))
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *GardenHandler) AdminCreateSection(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Slug string `json:"slug" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sec, err := h.garden.CreateSection(req.Name, req.Slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sec)
}

func (h *GardenHandler) AdminUpdateSection(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Name      string `json:"name"`
		SortOrder int    `json:"sort_order"`
		IsVisible bool   `json:"is_visible"`
	}
	c.ShouldBindJSON(&req)
	h.garden.UpdateSection(uint(id), req.Name, req.SortOrder, req.IsVisible)
	c.Status(http.StatusNoContent)
}

func (h *GardenHandler) AdminDeleteSection(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.garden.DeleteSection(uint(id))
	c.Status(http.StatusNoContent)
}

func (h *GardenHandler) AdminCreateItem(c *gin.Context) {
	var req struct {
		SectionID uint   `json:"section_id" binding:"required"`
		Type      string `json:"type" binding:"required"`
		Title     string `json:"title" binding:"required"`
		URL       string `json:"url"`
		Source    string `json:"source"`
		Summary   string `json:"summary"`
		EmbedURL  string `json:"embed_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.garden.CreateItem(req.SectionID, req.Type, req.Title, req.URL, req.Source, req.Summary, req.EmbedURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *GardenHandler) AdminUpdateItem(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Title     string `json:"title"`
		URL       string `json:"url"`
		Source    string `json:"source"`
		Summary   string `json:"summary"`
		EmbedURL  string `json:"embed_url"`
		SortOrder int    `json:"sort_order"`
	}
	c.ShouldBindJSON(&req)
	h.garden.UpdateItem(uint(id), req.Title, req.URL, req.Source, req.Summary, req.EmbedURL, req.SortOrder)
	c.Status(http.StatusNoContent)
}

func (h *GardenHandler) AdminDeleteItem(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.garden.DeleteItem(uint(id))
	c.Status(http.StatusNoContent)
}
```

Create `backend/internal/handler/stats.go`:

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/daiyutong/blog/internal/service"
	"github.com/gin-gonic/gin"
)

type StatsHandler struct{ stats *service.StatsService }

func NewStatsHandler(s *service.StatsService) *StatsHandler { return &StatsHandler{stats: s} }

func (h *StatsHandler) TrackPage(c *gin.Context) {
	var req struct {
		Path     string `json:"path"`
		Referrer string `json:"referrer"`
	}
	c.ShouldBindJSON(&req)
	hash := service.VisitorHash(c.ClientIP(), c.GetHeader("User-Agent"))
	h.stats.TrackPage(req.Path, hash, req.Referrer)
	c.Status(http.StatusNoContent)
}

func (h *StatsHandler) TrackPost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	hash := service.VisitorHash(c.ClientIP(), c.GetHeader("User-Agent"))
	h.stats.TrackPost(uint(id), hash)
	c.Status(http.StatusNoContent)
}

func (h *StatsHandler) AdminStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	c.JSON(http.StatusOK, h.stats.AdminStats(days))
}
```

Create `backend/internal/handler/upload.go`:

```go
package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{ uploadDir string }

func NewUploadHandler(uploadDir string) *UploadHandler {
	os.MkdirAll(uploadDir, 0755)
	return &UploadHandler{uploadDir: uploadDir}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file"})
		return
	}
	ext := filepath.Ext(file.Filename)
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(h.uploadDir, name)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": "/uploads/" + name})
}
```

- [ ] **Step 5: Commit**

```bash
git add backend/internal/handler/
git commit -m "feat: all HTTP handlers"
```

---

## Task 11: Router + Main

**Files:**
- Create: `backend/internal/router/router.go`
- Create: `backend/cmd/server/main.go`

- [ ] **Step 1: Write router**

Create `backend/internal/router/router.go`:

```go
package router

import (
	"github.com/daiyutong/blog/internal/handler"
	"github.com/daiyutong/blog/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth    *handler.AuthHandler
	Post    *handler.PostHandler
	Column  *handler.ColumnHandler
	Project *handler.ProjectHandler
	Garden  *handler.GardenHandler
	Stats   *handler.StatsHandler
	Upload  *handler.UploadHandler
}

func Setup(h Handlers, jwtSecret, uploadDir string) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	// Serve uploaded files
	r.Static("/uploads", uploadDir)

	api := r.Group("/api")

	// Public
	api.POST("/auth/login", h.Auth.Login)
	api.GET("/posts", h.Post.List)
	api.GET("/posts/:slug", h.Post.GetBySlug)
	api.GET("/tags", h.Post.Tags)
	api.GET("/columns", h.Column.List)
	api.GET("/columns/:slug", h.Column.GetBySlug)
	api.GET("/projects", h.Project.List)
	api.GET("/garden/sections", h.Garden.ListSections)
	api.GET("/garden/items", h.Garden.ListItems)
	api.POST("/track/page", h.Stats.TrackPage)
	api.POST("/track/post/:id", h.Stats.TrackPost)

	// Admin (JWT protected)
	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(jwtSecret))

	admin.GET("/stats", h.Stats.AdminStats)
	admin.POST("/posts", h.Post.AdminCreate)
	admin.PUT("/posts/:id", h.Post.AdminUpdate)
	admin.DELETE("/posts/:id", h.Post.AdminDelete)
	admin.PUT("/posts/:id/pin", h.Post.AdminSetPin)
	admin.PUT("/posts/:id/publish", h.Post.AdminSetPublished)
	admin.POST("/columns", h.Column.AdminCreate)
	admin.PUT("/columns/:id", h.Column.AdminUpdate)
	admin.DELETE("/columns/:id", h.Column.AdminDelete)
	admin.PUT("/columns/:id/posts", h.Column.AdminSetPosts)
	admin.POST("/projects", h.Project.AdminCreate)
	admin.PUT("/projects/:id", h.Project.AdminUpdate)
	admin.DELETE("/projects/:id", h.Project.AdminDelete)
	admin.POST("/garden/sections", h.Garden.AdminCreateSection)
	admin.PUT("/garden/sections/:id", h.Garden.AdminUpdateSection)
	admin.DELETE("/garden/sections/:id", h.Garden.AdminDeleteSection)
	admin.POST("/garden/items", h.Garden.AdminCreateItem)
	admin.PUT("/garden/items/:id", h.Garden.AdminUpdateItem)
	admin.DELETE("/garden/items/:id", h.Garden.AdminDeleteItem)
	admin.POST("/upload", h.Upload.Upload)

	return r
}
```

- [ ] **Step 2: Write main.go**

Create `backend/cmd/server/main.go`:

```go
package main

import (
	"log"

	"github.com/daiyutong/blog/internal/config"
	"github.com/daiyutong/blog/internal/db"
	"github.com/daiyutong/blog/internal/handler"
	"github.com/daiyutong/blog/internal/router"
	"github.com/daiyutong/blog/internal/service"
)

func main() {
	cfg := config.Load()
	database := db.Init(cfg.DBPath)

	uploadDir := "./data/uploads"

	postSvc := service.NewPostService(database)
	colSvc := service.NewColumnService(database)
	projSvc := service.NewProjectService(database)
	gardenSvc := service.NewGardenService(database)
	statsSvc := service.NewStatsService(database)

	h := router.Handlers{
		Auth:    handler.NewAuthHandler(cfg.AdminUser, cfg.AdminPass, cfg.JWTSecret),
		Post:    handler.NewPostHandler(postSvc, statsSvc),
		Column:  handler.NewColumnHandler(colSvc),
		Project: handler.NewProjectHandler(projSvc),
		Garden:  handler.NewGardenHandler(gardenSvc),
		Stats:   handler.NewStatsHandler(statsSvc),
		Upload:  handler.NewUploadHandler(uploadDir),
	}

	r := router.Setup(h, cfg.JWTSecret, uploadDir)
	log.Printf("server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
```

- [ ] **Step 3: Build to verify compilation**

```bash
cd /Users/daiyutong/IdeaProjects/blog/backend
go build ./...
```

Expected: no errors

- [ ] **Step 4: Run all tests**

```bash
go test ./... -v
```

Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add backend/internal/router/ backend/cmd/
git commit -m "feat: router setup and main entry point"
```

---

## Task 12: Dockerfile + Smoke Test

**Files:**
- Create: `backend/Dockerfile`

- [ ] **Step 1: Write Dockerfile**

Create `backend/Dockerfile`:

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache sqlite-libs ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

- [ ] **Step 2: Build Docker image**

```bash
cd /Users/daiyutong/IdeaProjects/blog/backend
docker build -t blog-api:dev .
```

Expected: build succeeds, image size under 30MB

- [ ] **Step 3: Run smoke test**

```bash
docker run -d --name blog-api-test -p 8080:8080 blog-api:dev
sleep 2
curl -s http://localhost:8080/api/posts | python3 -m json.tool
curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | python3 -m json.tool
docker rm -f blog-api-test
```

Expected: `/api/posts` returns `{"posts":[],"total":0,"page":1}`, login returns `{"token":"..."}` 

- [ ] **Step 4: Commit**

```bash
git add backend/Dockerfile
git commit -m "feat: backend Dockerfile, smoke test passed"
```

---

## Self-Review

**Spec coverage check:**
- ✅ All public API routes covered (posts, columns, projects, garden, track)
- ✅ All admin routes covered (CRUD for all entities + stats)
- ✅ JWT auth middleware on all admin routes
- ✅ Word count (goldmark) + reading minutes + view count on post detail
- ✅ Pinned posts (is_pinned sorted first in List)
- ✅ Column ordered posts (column_posts.sort_order)
- ✅ Visitor hash anonymization (SHA256, no raw IP stored)
- ✅ 24h dedup for post views
- ✅ daily_stats model defined (aggregation can be added as a cron goroutine in Phase 1.5)
- ✅ Image upload endpoint
- ✅ Docker single-binary build
- ✅ sort_order fields on projects, garden_sections, garden_items, column_posts (Phase 2 prep)

**Placeholder scan:** None found.

**Type consistency:**
- `service.WordCount` used in `markdown.go` and called in `post.go` ✅
- `service.ReadingMinutes` used in `markdown.go` and called in `handler/post.go` ✅
- `service.VisitorHash` used in `stats.go` and called in `handler/stats.go` ✅
- All model field names consistent between service and handler layers ✅
