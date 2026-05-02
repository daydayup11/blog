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

func (h *ColumnHandler) AdminList(c *gin.Context) {
	cols, err := h.columns.ListAll()
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

func (h *ColumnHandler) AdminGetPosts(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	posts, err := h.columns.GetPosts(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
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
		Published   *bool  `json:"published"`
	}
	c.ShouldBindJSON(&req)
	if req.Published != nil {
		h.columns.SetPublished(uint(id), *req.Published)
	}
	if req.Name != "" {
		h.columns.Update(uint(id), req.Name, req.Description, req.CoverURL)
	}
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
