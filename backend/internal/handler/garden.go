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

func (h *GardenHandler) ListWorksSections(c *gin.Context) {
	sections, _ := h.garden.ListWorksSections()
	c.JSON(http.StatusOK, gin.H{"sections": sections})
}

func (h *GardenHandler) AdminListAllSections(c *gin.Context) {
	sections, _ := h.garden.ListAllSections()
	c.JSON(http.StatusOK, gin.H{"sections": sections})
}

func (h *GardenHandler) AdminCreateWorksSection(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sec, err := h.garden.CreateWorksSection(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sec)
}

func (h *GardenHandler) ListItems(c *gin.Context) {
	sectionID, _ := strconv.ParseUint(c.Query("section_id"), 10, 64)
	items, _ := h.garden.ListItems(uint(sectionID), c.Query("type"), c.Query("tag"))
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
		Tags      string `json:"tags"`
		EmbedURL  string `json:"embed_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Tags == "" {
		req.Tags = "[]"
	}
	item, err := h.garden.CreateItem(req.SectionID, req.Type, req.Title, req.URL, req.Source, req.Summary, req.Tags, req.EmbedURL)
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
		Tags      string `json:"tags"`
		EmbedURL  string `json:"embed_url"`
		SortOrder int    `json:"sort_order"`
	}
	c.ShouldBindJSON(&req)
	h.garden.UpdateItem(uint(id), req.Title, req.URL, req.Source, req.Summary, req.Tags, req.EmbedURL, req.SortOrder)
	c.Status(http.StatusNoContent)
}

func (h *GardenHandler) AdminDeleteItem(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.garden.DeleteItem(uint(id))
	c.Status(http.StatusNoContent)
}
