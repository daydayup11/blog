package handler

import (
	"net/http"
	"strconv"

	"github.com/daiyutong/blog/internal/model"
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
	posts, total, err := h.posts.List(page, 10, c.Query("tag"), c.Query("q"))
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
	prev, next := h.posts.GetAdjacentPosts(post.CreatedAt)

	toSlim := func(p *model.Post) interface{} {
		if p == nil {
			return nil
		}
		return gin.H{"title": p.Title, "slug": p.Slug}
	}

	c.JSON(http.StatusOK, gin.H{
		"post":            post,
		"views":           views,
		"reading_minutes": service.ReadingMinutes(post.WordCount),
		"prev":            toSlim(prev),
		"next":            toSlim(next),
	})
}

func (h *PostHandler) Tags(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tags": h.posts.AllTags()})
}

func (h *PostHandler) AdminGetBySlug(c *gin.Context) {
	post, err := h.posts.GetBySlugAdmin(c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	views := h.stats.PostViewCount(post.ID)
	// Content has json:"-" to hide from public API; expose it here for the CMS editor
	c.JSON(http.StatusOK, gin.H{
		"post": gin.H{
			"id":           post.ID,
			"title":        post.Title,
			"slug":         post.Slug,
			"content":      post.Content,
			"summary":      post.Summary,
			"word_count":   post.WordCount,
			"tags":         post.Tags,
			"is_pinned":    post.IsPinned,
			"is_published": post.IsPublished,
			"created_at":   post.CreatedAt,
			"updated_at":   post.UpdatedAt,
		},
		"views":           views,
		"reading_minutes": service.ReadingMinutes(post.WordCount),
	})
}

func (h *PostHandler) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	posts, total, err := h.posts.ListAll(page, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts, "total": total, "page": page})
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
