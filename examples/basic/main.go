package main

import (
	"fmt"
	"log"

	"github.com/nate3d/toolbox/pkg/utils"
)

// Example demonstrates basic usage of the toolbox utilities
func main() {
	fmt.Println("Go Toolbox Examples")
	fmt.Println("===================")

	// String utilities
	fmt.Println("\n1. String Utilities:")
	str := utils.String()

	text := "hello world"
	fmt.Printf("Original: %s\n", text)
	fmt.Printf("Reversed: %s\n", str.Reverse(text))
	fmt.Printf("CamelCase: %s\n", str.ToCamelCase(text))
	fmt.Printf("SnakeCase: %s\n", str.ToSnakeCase("HelloWorld"))
	fmt.Printf("KebabCase: %s\n", str.ToKebabCase("HelloWorld"))

	// Slice utilities
	fmt.Println("\n2. Slice Utilities:")
	slice := utils.Slice()

	fruits := []string{"apple", "banana", "apple", "cherry", "banana"}
	fmt.Printf("Original: %v\n", fruits)
	fmt.Printf("Unique: %v\n", slice.Unique(fruits))
	fmt.Printf("Contains 'apple': %v\n", slice.Contains(fruits, "apple"))
	fmt.Printf("Sorted: %v\n", slice.Sort(fruits))

	// Filter example
	longFruits := slice.Filter(fruits, func(s string) bool {
		return len(s) > 5
	})
	fmt.Printf("Long fruits (>5 chars): %v\n", longFruits)

	// Random utilities
	fmt.Println("\n3. Random Utilities:")
	random := utils.Random()

	fmt.Printf("Random string (10 chars): %s\n", random.String(10))
	fmt.Printf("Random int (1-100): %d\n", random.Int(1, 100))
	fmt.Printf("Random bool: %v\n", random.Bool())
	fmt.Printf("Random choice from fruits: %s\n", random.Choice([]string{"apple", "banana", "cherry"}))

	// Hash utilities
	fmt.Println("\n4. Hash Utilities:")
	hash := utils.Hash()

	data := "hello world"
	fmt.Printf("Text: %s\n", data)
	fmt.Printf("MD5: %s\n", hash.MD5(data))
	fmt.Printf("SHA256: %s\n", hash.SHA256(data))

	// Validation utilities
	fmt.Println("\n5. Validation Utilities:")
	validate := utils.Validate()

	emails := []string{"test@example.com", "invalid.email", "user@domain.co.uk"}
	for _, email := range emails {
		fmt.Printf("Email '%s' is valid: %v\n", email, validate.Email(email))
	}

	ips := []string{"192.168.1.1", "10.0.0.1", "256.1.1.1", "not.an.ip"}
	for _, ip := range ips {
		fmt.Printf("IP '%s' is valid: %v\n", ip, validate.IP(ip))
	}

	// File utilities
	fmt.Println("\n6. File Utilities:")
	file := utils.File()

	// Check if files exist
	testFiles := []string{"go.mod", "nonexistent.txt", "README.md"}
	for _, filename := range testFiles {
		fmt.Printf("File '%s' exists: %v\n", filename, file.Exists(filename))
	}

	// Get file size
	if file.Exists("go.mod") {
		if size, err := file.Size("go.mod"); err == nil {
			fmt.Printf("go.mod size: %d bytes\n", size)
		}
	}

	// Conversion utilities
	fmt.Println("\n7. Conversion Utilities:")
	convert := utils.Convert()

	// String conversions
	if num, err := convert.StringToInt("42"); err == nil {
		fmt.Printf("String '42' to int: %d\n", num)
	}

	if f, err := convert.StringToFloat("3.14"); err == nil {
		fmt.Printf("String '3.14' to float: %.2f\n", f)
	}

	if b, err := convert.StringToBool("true"); err == nil {
		fmt.Printf("String 'true' to bool: %v\n", b)
	}

	// Convert back to strings
	fmt.Printf("Int 42 to string: '%s'\n", convert.IntToString(42))
	fmt.Printf("Float 3.14 to string: '%s'\n", convert.FloatToString(3.14))
	fmt.Printf("Bool true to string: '%s'\n", convert.BoolToString(true))

	fmt.Println("\nExample completed successfully!")
}

// Advanced examples

// ExampleConcurrency demonstrates how you might use these utilities in a concurrent context
func ExampleConcurrency() {
	fmt.Println("\nConcurrency Example:")

	// Process multiple URLs concurrently (simulated)
	urls := []string{
		"https://example.com",
		"https://google.com",
		"https://github.com",
	}

	hash := utils.Hash()

	for _, url := range urls {
		go func(u string) {
			// In a real scenario, you'd fetch the URL content
			// Here we just hash the URL itself
			urlHash := hash.MD5(u)
			fmt.Printf("URL: %s, Hash: %s\n", u, urlHash)
		}(url)
	}
}

// ExampleChaining demonstrates chaining utility operations
func ExampleChaining() {
	fmt.Println("\nChaining Example:")

	str := utils.String()
	slice := utils.Slice()

	// Process a list of mixed-case strings
	names := []string{"john DOE", "jane_smith", "bob-jones", "alice WILLIAMS"}

	// Convert all to camelCase, make unique, and sort
	processed := make([]string, len(names))
	for i, name := range names {
		processed[i] = str.ToCamelCase(name)
	}

	processed = slice.Unique(processed)
	processed = slice.Sort(processed)

	fmt.Printf("Original: %v\n", names)
	fmt.Printf("Processed: %v\n", processed)
}

// ExampleErrorHandling demonstrates proper error handling with utilities
func ExampleErrorHandling() {
	fmt.Println("\nError Handling Example:")

	convert := utils.Convert()
	validate := utils.Validate()

	// Try to convert invalid strings
	invalidInputs := []string{"not_a_number", "3.14.15", "maybe"}

	for _, input := range invalidInputs {
		if num, err := convert.StringToInt(input); err != nil {
			fmt.Printf("Failed to convert '%s' to int: %v\n", input, err)
		} else {
			fmt.Printf("Converted '%s' to int: %d\n", input, num)
		}
	}

	// Validate emails with error context
	emails := []string{"good@example.com", "bad.email", "also@good.com"}

	for _, email := range emails {
		if validate.Email(email) {
			fmt.Printf("✓ Valid email: %s\n", email)
		} else {
			fmt.Printf("✗ Invalid email: %s\n", email)
		}
	}
}

// init function runs before main
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
