package service_test

import (
	"testing"

	"github.com/daiyutong/blog/internal/service"
)

func TestVisitorHash_Deterministic(t *testing.T) {
	h1 := service.VisitorHash("192.168.1.1", "Mozilla/5.0")
	h2 := service.VisitorHash("192.168.1.1", "Mozilla/5.0")
	if h1 != h2 {
		t.Error("same inputs should produce same hash")
	}
}

func TestVisitorHash_Different(t *testing.T) {
	h1 := service.VisitorHash("192.168.1.1", "Mozilla/5.0")
	h2 := service.VisitorHash("10.0.0.1", "Mozilla/5.0")
	if h1 == h2 {
		t.Error("different IPs should produce different hashes")
	}
}

func TestVisitorHash_NotRawIP(t *testing.T) {
	h := service.VisitorHash("192.168.1.1", "agent")
	if h == "192.168.1.1" {
		t.Error("hash must not be the raw IP")
	}
	if len(h) < 16 {
		t.Error("hash too short")
	}
}
