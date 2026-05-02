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
