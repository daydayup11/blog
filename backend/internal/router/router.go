package router

import (
	"net/http"
	"path/filepath"
	"strings"

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

func Setup(h Handlers, jwtSecret, uploadDir, frontendDir string) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	r.Static("/uploads", uploadDir)

	// Serve frontend static files
	r.Static("/css", filepath.Join(frontendDir, "css"))
	r.Static("/js",  filepath.Join(frontendDir, "js"))
	r.Static("/admin", filepath.Join(frontendDir, "admin"))
	// Serve root HTML pages
	for _, page := range []string{"index", "blog", "post", "column", "works", "garden"} {
		p := page // capture
		r.GET("/"+p+".html", func(c *gin.Context) {
			c.File(filepath.Join(frontendDir, p+".html"))
		})
	}
	r.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(frontendDir, "index.html"))
	})
	// Fallback: serve index.html for unknown paths (not /api, not /uploads)
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/uploads/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.File(filepath.Join(frontendDir, "index.html"))
	})

	api := r.Group("/api")

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

	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(jwtSecret))

	admin.GET("/stats", h.Stats.AdminStats)
	admin.GET("/posts", h.Post.AdminList)
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
