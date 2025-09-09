// Package logger provides structured logging utilities for the toolbox applications.
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger

	level  slog.Level
	output io.Writer
}

// LogLevel represents the logging level
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// Config holds logger configuration
type Config struct {
	Level      LogLevel
	Output     string // "stdout", "stderr", or file path
	Format     string // "text" or "json"
	WithCaller bool
	WithTime   bool
}

var globalLogger *Logger

// Init initializes the global logger
func Init(config Config) error {
	level := parseLevel(config.Level)

	// Determine output writer
	var output io.Writer
	switch config.Output {
	case "", "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		// File output
		if err := os.MkdirAll(filepath.Dir(config.Output), 0750); err != nil {
			return err
		}
		file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
		output = file
	}

	// Create handler based on format
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: config.WithCaller,
	}

	if config.Format == "json" {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = NewColorHandler(output, opts)
	}

	logger := slog.New(handler)
	globalLogger = &Logger{
		Logger: logger,
		level:  level,
		output: output,
	}

	return nil
}

// Get returns the global logger
func Get() *Logger {
	if globalLogger == nil {
		// Initialize with default config if not initialized
		_ = Init(Config{
			Level:      LevelInfo,
			Output:     "stdout",
			Format:     "text",
			WithCaller: false,
			WithTime:   true,
		})
	}
	return globalLogger
}

// parseLevel converts string level to slog.Level
func parseLevel(level LogLevel) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// ColorHandler is a custom slog handler that adds colors to console output
type ColorHandler struct {
	handler slog.Handler
	output  io.Writer
}

// NewColorHandler creates a new ColorHandler
func NewColorHandler(output io.Writer, opts *slog.HandlerOptions) *ColorHandler {
	return &ColorHandler{
		handler: slog.NewTextHandler(output, opts),
		output:  output,
	}
}

// Enabled implements slog.Handler
func (h *ColorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle implements slog.Handler
func (h *ColorHandler) Handle(ctx context.Context, record slog.Record) error {
	// Check if output supports color (is a terminal)
	if file, ok := h.output.(*os.File); ok && isTerminal(file) {
		return h.handleWithColor(ctx, record)
	}
	return h.handler.Handle(ctx, record)
}

// WithAttrs implements slog.Handler
func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ColorHandler{
		handler: h.handler.WithAttrs(attrs),
		output:  h.output,
	}
}

// WithGroup implements slog.Handler
func (h *ColorHandler) WithGroup(name string) slog.Handler {
	return &ColorHandler{
		handler: h.handler.WithGroup(name),
		output:  h.output,
	}
}

// handleWithColor handles the record with color formatting.
func (h *ColorHandler) handleWithColor(_ context.Context, record slog.Record) error {
	var levelColor *color.Color
	var levelText string

	switch record.Level {
	case slog.LevelDebug:
		levelColor = color.New(color.FgCyan)
		levelText = "DEBUG"
	case slog.LevelInfo:
		levelColor = color.New(color.FgGreen)
		levelText = "INFO "
	case slog.LevelWarn:
		levelColor = color.New(color.FgYellow)
		levelText = "WARN "
	case slog.LevelError:
		levelColor = color.New(color.FgRed)
		levelText = "ERROR"
	default:
		levelColor = color.New(color.Reset)
		levelText = "UNKNOWN"
	}

	// Format timestamp
	timestamp := record.Time.Format("2006-01-02 15:04:05")

	// Write colored output
	_, _ = color.New(color.FgHiBlack).Fprintf(h.output, "%s ", timestamp)
	_, _ = levelColor.Fprintf(h.output, "[%s] ", levelText)

	// Write message
	if record.Level >= slog.LevelError {
		_, _ = color.New(color.FgRed).Fprintf(h.output, "%s", record.Message)
	} else {
		_, _ = fmt.Fprintf(h.output, "%s", record.Message)
	}

	// Write attributes
	record.Attrs(func(attr slog.Attr) bool {
		_, _ = color.New(color.FgHiBlack).Fprintf(h.output, " %s=%v", attr.Key, attr.Value)
		return true
	})

	_, _ = fmt.Fprintln(h.output)
	return nil
}

// isTerminal checks if the file is a terminal
func isTerminal(file *os.File) bool {
	stat, err := file.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// Convenience methods for the global logger

// Debug logs a debug message
func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

// Info logs an info message
func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

// With returns a logger with the given attributes
func With(args ...any) *Logger {
	return &Logger{
		Logger: Get().With(args...),
		level:  Get().level,
		output: Get().output,
	}
}
