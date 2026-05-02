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

	postSvc    := service.NewPostService(database)
	profileSvc := service.NewProfileService(database)
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
		Profile: handler.NewProfileHandler(profileSvc),
	}

	r := router.Setup(h, cfg.JWTSecret, uploadDir, cfg.FrontendDir)
	log.Printf("server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
