// Package config provides configuration management utilities for the toolbox applications.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the global configuration structure
type Config struct {
	// Global settings
	LogLevel string `mapstructure:"log_level"`
	LogFile  string `mapstructure:"log_file"`

	// Application-specific settings
	CLI CLIConfig `mapstructure:"cli"`
	TUI TUIConfig `mapstructure:"tui"`
	Web WebConfig `mapstructure:"web"`
}

// CLIConfig holds CLI-specific configuration
type CLIConfig struct {
	DefaultOutput string `mapstructure:"default_output"`
	ColorOutput   bool   `mapstructure:"color_output"`
	Verbose       bool   `mapstructure:"verbose"`
}

// TUIConfig holds TUI-specific configuration
type TUIConfig struct {
	Theme       string `mapstructure:"theme"`
	MouseEvents bool   `mapstructure:"mouse_events"`
}

// WebConfig holds web-specific configuration
type WebConfig struct {
	Port    int    `mapstructure:"port"`
	Host    string `mapstructure:"host"`
	TLSCert string `mapstructure:"tls_cert"`
	TLSKey  string `mapstructure:"tls_key"`
}

var globalConfig *Config

// Init initializes the configuration system
func Init(appName string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Add configuration paths
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("$HOME/.config/" + appName)
	viper.AddConfigPath("/etc/" + appName)
	viper.AddConfigPath(".")

	// Set environment variable prefix
	viper.SetEnvPrefix(strings.ToUpper(appName))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		var configNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configNotFoundError) {
			return fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, use defaults
	}

	// Unmarshal into struct
	globalConfig = &Config{}
	if err := viper.Unmarshal(globalConfig); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Global defaults
	viper.SetDefault("log_level", "info")
	viper.SetDefault("log_file", "")

	// CLI defaults
	viper.SetDefault("cli.default_output", "table")
	viper.SetDefault("cli.color_output", true)
	viper.SetDefault("cli.verbose", false)

	// TUI defaults
	viper.SetDefault("tui.theme", "default")
	viper.SetDefault("tui.mouse_events", true)

	// Web defaults
	viper.SetDefault("web.port", 8080)
	viper.SetDefault("web.host", "localhost")
	viper.SetDefault("web.tls_cert", "")
	viper.SetDefault("web.tls_key", "")
}

// Get returns the global configuration
func Get() *Config {
	if globalConfig == nil {
		// Initialize with default values if not initialized
		globalConfig = &Config{}
		setDefaults()
		_ = viper.Unmarshal(globalConfig)
	}
	return globalConfig
}

// GetString returns a configuration value as string
func GetString(key string) string {
	return viper.GetString(key)
}

// GetBool returns a configuration value as bool
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetInt returns a configuration value as int
func GetInt(key string) int {
	return viper.GetInt(key)
}

// Set sets a configuration value
func Set(key string, value interface{}) {
	viper.Set(key, value)
}

// WriteConfig writes the current configuration to file
func WriteConfig() error {
	return viper.WriteConfig()
}

// WriteConfigAs writes the current configuration to a specific file
func WriteConfigAs(filename string) error {
	return viper.WriteConfigAs(filename)
}

// GetConfigDir returns the configuration directory for the application
func GetConfigDir(appName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", appName)
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return "", err
	}

	return configDir, nil
}
