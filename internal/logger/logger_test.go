package logger

import (
	"testing"
)

func TestInitAndGet(t *testing.T) {
	cfg := Config{Level: LevelDebug, Output: "stdout", Format: "text", WithCaller: false, WithTime: false}
	if err := Init(cfg); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	log := Get()
	if log == nil {
		t.Fatal("Get returned nil")
	}
}
