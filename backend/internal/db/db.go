package db

import (
	"log"
	"os"
	"path/filepath"

	"github.com/daiyutong/blog/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(dbPath string) *gorm.DB {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("failed to create db dir: %v", err)
	}
	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	if err := database.AutoMigrate(
		&model.Post{},
		&model.Column{},
		&model.ColumnPost{},
		&model.Project{},
		&model.GardenSection{},
		&model.GardenItem{},
		&model.PostView{},
		&model.PageView{},
		&model.DailyStat{},
		&model.Profile{},
	); err != nil {
		log.Fatalf("automigrate failed: %v", err)
	}
	return database
}
