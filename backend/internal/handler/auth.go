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
