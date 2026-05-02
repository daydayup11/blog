package handler

import (
	"net/http"

	"github.com/daiyutong/blog/internal/service"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profile *service.ProfileService
}

func NewProfileHandler(p *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profile: p}
}

func (h *ProfileHandler) Get(c *gin.Context) {
	p, err := h.profile.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProfileHandler) Save(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Subtitle string `json:"subtitle"`
		Bio      string `json:"bio"`
		Skills   string `json:"skills"`
		Social   string `json:"social"`
		Timeline string `json:"timeline"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Skills == "" {
		req.Skills = "[]"
	}
	if req.Social == "" {
		req.Social = "[]"
	}
	if req.Timeline == "" {
		req.Timeline = "[]"
	}
	p, err := h.profile.Save(req.Name, req.Subtitle, req.Bio, req.Skills, req.Social, req.Timeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}
