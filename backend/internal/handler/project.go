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
