package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigValues(t *testing.T) {
	// Initialize with a dummy app name, config file may not exist
	if err := Init("testapp"); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	cfg := Get()
	// Check global defaults
	if cfg.LogLevel != "info" {
		t.Errorf("Default LogLevel = %q, want %q", cfg.LogLevel, "info")
	}
	if cfg.LogFile != "" {
		t.Errorf("Default LogFile = %q, want empty", cfg.LogFile)
	}
	// Check CLI defaults
	if cfg.CLI.DefaultOutput != "table" {
		t.Errorf("CLI.DefaultOutput = %q, want %q", cfg.CLI.DefaultOutput, "table")
	}
	if cfg.CLI.ColorOutput != true {
		t.Errorf("CLI.ColorOutput = %v, want %v", cfg.CLI.ColorOutput, true)
	}
	if cfg.CLI.Verbose != false {
		t.Errorf("CLI.Verbose = %v, want %v", cfg.CLI.Verbose, false)
	}
	// Check Web defaults
	if cfg.Web.Port != 8080 {
		t.Errorf("Web.Port = %d, want %d", cfg.Web.Port, 8080)
	}
	if cfg.Web.Host != "localhost" {
		t.Errorf("Web.Host = %q, want %q", cfg.Web.Host, "localhost")
	}
}

func TestGetConfigDirCreatesDirectory(t *testing.T) {
	appName := "testapp"
	// Clean up any previous state
	d := filepath.Join(os.TempDir(), ".config", appName)
	os.RemoveAll(d)

	// Override UserHomeDir to temp dir by setting HOME env var
	t.Setenv("HOME", os.TempDir())

	dir, err := GetConfigDir(appName)
	if err != nil {
		t.Fatalf("GetConfigDir failed: %v", err)
	}
	// Should exist
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("Configured directory not found: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("Expected %s to be a directory", dir)
	}
}
