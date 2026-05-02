package service_test

import (
	"testing"

	"github.com/daiyutong/blog/internal/service"
)

func TestWordCount_ChineseAndEnglish(t *testing.T) {
	md := "Hello world 你好世界"
	count := service.WordCount(md)
	if count != 6 {
		t.Errorf("expected 6, got %d", count)
	}
}

func TestWordCount_StripsMarkdown(t *testing.T) {
	md := "# 标题\n\n**加粗** 文字"
	count := service.WordCount(md)
	if count != 6 {
		t.Errorf("expected 6, got %d", count)
	}
}

func TestRenderMarkdown_ReturnsHTML(t *testing.T) {
	md := "# Hello\n\nworld"
	html := service.RenderMarkdown(md)
	if html == "" {
		t.Error("expected non-empty HTML")
	}
	if len(html) < 10 {
		t.Errorf("HTML too short: %s", html)
	}
}

func TestReadingMinutes(t *testing.T) {
	if service.ReadingMinutes(0) != 1 {
		t.Error("0 words should return 1 min minimum")
	}
	if service.ReadingMinutes(300) != 1 {
		t.Error("300 words should return 1 min")
	}
	if service.ReadingMinutes(600) != 2 {
		t.Error("600 words should return 2 min")
	}
}
