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
