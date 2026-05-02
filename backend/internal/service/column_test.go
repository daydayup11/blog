package service_test

import (
	"testing"

	"github.com/daiyutong/blog/internal/db"
	"github.com/daiyutong/blog/internal/service"
)

func setupColumnDB(t *testing.T) (*service.ColumnService, *service.PostService) {
	t.Helper()
	database := db.Init(":memory:")
	return service.NewColumnService(database), service.NewPostService(database)
}

func TestCreateColumn(t *testing.T) {
	svc, _ := setupColumnDB(t)
	col, err := svc.Create("Go 系列", "go-series", "Go 从入门到实践", "")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if col.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestColumnPostOrder(t *testing.T) {
	colSvc, postSvc := setupColumnDB(t)
	col, _ := colSvc.Create("系列", "series", "", "")
	p1, _ := postSvc.Create("第一篇", "post-1", "内容", "", []string{})
	p2, _ := postSvc.Create("第二篇", "post-2", "内容", "", []string{})

	err := colSvc.SetPosts(col.ID, []uint{p1.ID, p2.ID})
	if err != nil {
		t.Fatalf("set posts failed: %v", err)
	}

	posts, err := colSvc.GetPosts(col.ID)
	if err != nil {
		t.Fatalf("get posts failed: %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
	if posts[0].ID != p1.ID {
		t.Error("first post should be p1")
	}
}
