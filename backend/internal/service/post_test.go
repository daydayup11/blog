package service_test

import (
	"testing"

	"github.com/daiyutong/blog/internal/db"
	"github.com/daiyutong/blog/internal/service"
)

func setupTestDB(t *testing.T) *service.PostService {
	t.Helper()
	database := db.Init(":memory:")
	return service.NewPostService(database)
}

func TestCreatePost_SetsWordCount(t *testing.T) {
	svc := setupTestDB(t)
	post, err := svc.Create("测试标题", "test-slug", "你好世界 hello", "摘要", []string{"go"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if post.WordCount != 5 {
		t.Errorf("expected word count 5, got %d", post.WordCount)
	}
}

func TestCreatePost_SlugUnique(t *testing.T) {
	svc := setupTestDB(t)
	_, err := svc.Create("标题1", "same-slug", "内容", "", []string{})
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}
	_, err = svc.Create("标题2", "same-slug", "内容2", "", []string{})
	if err == nil {
		t.Error("expected error for duplicate slug")
	}
}

func TestSetPinned(t *testing.T) {
	svc := setupTestDB(t)
	post, _ := svc.Create("标题", "pin-slug", "内容", "", []string{})
	if err := svc.SetPinned(post.ID, true); err != nil {
		t.Fatalf("set pinned failed: %v", err)
	}
	svc.SetPublished(post.ID, true)
	fetched, _ := svc.GetBySlug("pin-slug")
	if !fetched.IsPinned {
		t.Error("expected post to be pinned")
	}
}

func TestListPosts_PinnedFirst(t *testing.T) {
	svc := setupTestDB(t)
	p1, _ := svc.Create("普通", "normal", "内容", "", []string{})
	p2, _ := svc.Create("置顶", "pinned", "内容", "", []string{})
	svc.SetPinned(p2.ID, true)
	svc.SetPublished(p2.ID, true)
	svc.SetPublished(p1.ID, true)

	posts, _, _ := svc.List(1, 10, "", "")
	if len(posts) == 0 {
		t.Fatal("expected posts")
	}
	if !posts[0].IsPinned {
		t.Error("first post should be pinned")
	}
}
