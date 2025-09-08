// Package utils provides general utility functions.
package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// StringUtils provides string manipulation utilities
type StringUtils struct{}

// String returns a new StringUtils instance
func String() *StringUtils {
	return &StringUtils{}
}

// IsEmpty checks if a string is empty or contains only whitespace
func (s *StringUtils) IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// Reverse reverses a string
func (s *StringUtils) Reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Truncate truncates a string to a maximum length
func (s *StringUtils) Truncate(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	if maxLen <= 3 {
		return str[:maxLen]
	}
	return str[:maxLen-3] + "..."
}

// PadLeft pads a string to the left with the specified character
func (s *StringUtils) PadLeft(str string, totalLen int, padChar rune) string {
	strLen := len([]rune(str))
	if strLen >= totalLen {
		return str
	}
	return strings.Repeat(string(padChar), totalLen-strLen) + str
}

// PadRight pads a string to the right with the specified character
func (s *StringUtils) PadRight(str string, totalLen int, padChar rune) string {
	strLen := len([]rune(str))
	if strLen >= totalLen {
		return str
	}
	return str + strings.Repeat(string(padChar), totalLen-strLen)
}

// ToCamelCase converts a string to camelCase
func (s *StringUtils) ToCamelCase(str string) string {
	words := strings.FieldsFunc(str, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})

	if len(words) == 0 {
		return ""
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		result += strings.Title(strings.ToLower(words[i]))
	}
	return result
}

// ToSnakeCase converts a string to snake_case
func (s *StringUtils) ToSnakeCase(str string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

// ToKebabCase converts a string to kebab-case
func (s *StringUtils) ToKebabCase(str string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	kebab := re.ReplaceAllString(str, "${1}-${2}")
	return strings.ToLower(kebab)
}

// SliceUtils provides slice manipulation utilities
type SliceUtils struct{}

// Slice returns a new SliceUtils instance
func Slice() *SliceUtils {
	return &SliceUtils{}
}

// Contains checks if a string slice contains a specific item
func (s *SliceUtils) Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Unique removes duplicate strings from a slice
func (s *SliceUtils) Unique(slice []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// Filter filters a string slice based on a predicate function
func (s *SliceUtils) Filter(slice []string, predicate func(string) bool) []string {
	result := make([]string, 0)
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map applies a function to each element of a string slice
func (s *SliceUtils) Map(slice []string, mapper func(string) string) []string {
	result := make([]string, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}

// Chunk splits a slice into smaller chunks of specified size
func (s *SliceUtils) Chunk(slice []string, size int) [][]string {
	if size <= 0 {
		return nil
	}

	chunks := make([][]string, 0)
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// Sort sorts a string slice and returns a new slice
func (s *SliceUtils) Sort(slice []string) []string {
	result := make([]string, len(slice))
	copy(result, slice)
	sort.Strings(result)
	return result
}

// FileUtils provides file system utilities
type FileUtils struct{}

// File returns a new FileUtils instance
func File() *FileUtils {
	return &FileUtils{}
}

// Exists checks if a file or directory exists
func (f *FileUtils) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDir checks if a path is a directory
func (f *FileUtils) IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile checks if a path is a regular file
func (f *FileUtils) IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// Size returns the size of a file in bytes
func (f *FileUtils) Size(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// MkdirAll creates directories recursively
func (f *FileUtils) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// Copy copies a file from src to dst
func (f *FileUtils) Copy(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// ReadLines reads all lines from a file
func (f *FileUtils) ReadLines(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	// Remove last empty line if it exists
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines, nil
}

// WriteLines writes lines to a file
func (f *FileUtils) WriteLines(path string, lines []string) error {
	content := strings.Join(lines, "\n")
	return os.WriteFile(path, []byte(content), 0644)
}

// Glob returns all files matching a pattern
func (f *FileUtils) Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

// HashUtils provides hashing utilities
type HashUtils struct{}

// Hash returns a new HashUtils instance
func Hash() *HashUtils {
	return &HashUtils{}
}

// MD5 calculates the SHA256 hash of a string (MD5 is deprecated, using SHA256 instead)
func (h *HashUtils) MD5(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

// SHA256 calculates the SHA256 hash of a string
func (h *HashUtils) SHA256(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

// MD5File calculates the SHA256 hash of a file (MD5 is deprecated, using SHA256 instead)
func (h *HashUtils) MD5File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// SHA256File calculates the SHA256 hash of a file
func (h *HashUtils) SHA256File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// RandomUtils provides random generation utilities
type RandomUtils struct {
	rand *rand.Rand
}

// Random returns a new RandomUtils instance
func Random() *RandomUtils {
	return &RandomUtils{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// String generates a random string of specified length
func (r *RandomUtils) String(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.rand.Intn(len(charset))]
	}
	return string(result)
}

// Int generates a random integer between min and max (inclusive)
func (r *RandomUtils) Int(min, max int) int {
	return r.rand.Intn(max-min+1) + min
}

// Bool generates a random boolean
func (r *RandomUtils) Bool() bool {
	return r.rand.Intn(2) == 1
}

// Choice randomly selects an item from a slice
func (r *RandomUtils) Choice(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[r.rand.Intn(len(items))]
}

// Shuffle shuffles a string slice in place
func (r *RandomUtils) Shuffle(slice []string) {
	r.rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

// ValidationUtils provides validation utilities
type ValidationUtils struct{}

// Validate returns a new ValidationUtils instance
func Validate() *ValidationUtils {
	return &ValidationUtils{}
}

// Email validates an email address
func (v *ValidationUtils) Email(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// URL validates a URL
func (v *ValidationUtils) URL(url string) bool {
	re := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return re.MatchString(url)
}

// IP validates an IP address (IPv4)
func (v *ValidationUtils) IP(ip string) bool {
	re := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	if !re.MatchString(ip) {
		return false
	}

	parts := strings.Split(ip, ".")
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return false
		}
	}
	return true
}

// PhoneNumber validates a phone number (basic validation)
func (v *ValidationUtils) PhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^\+?[\d\s\-\(\)]+$`)
	return re.MatchString(phone) && len(phone) >= 10
}

// ConversionUtils provides conversion utilities
type ConversionUtils struct{}

// Convert returns a new ConversionUtils instance
func Convert() *ConversionUtils {
	return &ConversionUtils{}
}

// StringToInt converts a string to int with error handling
func (c *ConversionUtils) StringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

// StringToFloat converts a string to float64 with error handling
func (c *ConversionUtils) StringToFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

// StringToBool converts a string to bool with error handling
func (c *ConversionUtils) StringToBool(str string) (bool, error) {
	return strconv.ParseBool(str)
}

// IntToString converts an int to string
func (c *ConversionUtils) IntToString(i int) string {
	return strconv.Itoa(i)
}

// FloatToString converts a float64 to string
func (c *ConversionUtils) FloatToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

// BoolToString converts a bool to string
func (c *ConversionUtils) BoolToString(b bool) string {
	return strconv.FormatBool(b)
}
