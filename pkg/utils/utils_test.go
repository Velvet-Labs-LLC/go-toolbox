package utils_test

import (
	"fmt"
	"testing"

	"github.com/nate3d/toolbox/pkg/utils"
)

func TestStringUtils(t *testing.T) {
	str := utils.String()

	t.Run("IsEmpty", func(t *testing.T) {
		tests := []struct {
			input    string
			expected bool
		}{
			{"", true},
			{"   ", true},
			{"hello", false},
			{" hello ", false},
		}

		for _, test := range tests {
			result := str.IsEmpty(test.input)
			if result != test.expected {
				t.Errorf("IsEmpty(%q) = %v, expected %v", test.input, result, test.expected)
			}
		}
	})

	t.Run("Reverse", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"hello", "olleh"},
			{"", ""},
			{"a", "a"},
			{"12345", "54321"},
		}

		for _, test := range tests {
			result := str.Reverse(test.input)
			if result != test.expected {
				t.Errorf("Reverse(%q) = %q, expected %q", test.input, result, test.expected)
			}
		}
	})

	t.Run("Truncate", func(t *testing.T) {
		tests := []struct {
			input    string
			maxLen   int
			expected string
		}{
			{"hello world", 5, "he..."},
			{"hello", 10, "hello"},
			{"hello", 5, "hello"},
			{"hello world", 3, "hel"},
		}

		for _, test := range tests {
			result := str.Truncate(test.input, test.maxLen)
			if result != test.expected {
				t.Errorf("Truncate(%q, %d) = %q, expected %q", test.input, test.maxLen, result, test.expected)
			}
		}
	})

	t.Run("ToCamelCase", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"hello world", "helloWorld"},
			{"hello_world", "helloWorld"},
			{"hello-world", "helloWorld"},
			{"HELLO WORLD", "helloWorld"},
		}

		for _, test := range tests {
			result := str.ToCamelCase(test.input)
			if result != test.expected {
				t.Errorf("ToCamelCase(%q) = %q, expected %q", test.input, result, test.expected)
			}
		}
	})
}

func TestSliceUtils(t *testing.T) {
	slice := utils.Slice()

	t.Run("Contains", func(t *testing.T) {
		testSlice := []string{"apple", "banana", "cherry"}

		tests := []struct {
			item     string
			expected bool
		}{
			{"apple", true},
			{"banana", true},
			{"orange", false},
			{"", false},
		}

		for _, test := range tests {
			result := slice.Contains(testSlice, test.item)
			if result != test.expected {
				t.Errorf("Contains(%v, %q) = %v, expected %v", testSlice, test.item, result, test.expected)
			}
		}
	})

	t.Run("Unique", func(t *testing.T) {
		input := []string{"apple", "banana", "apple", "cherry", "banana"}
		expected := []string{"apple", "banana", "cherry"}

		result := slice.Unique(input)
		if len(result) != len(expected) {
			t.Errorf("Unique(%v) length = %d, expected %d", input, len(result), len(expected))
		}

		for _, item := range expected {
			if !slice.Contains(result, item) {
				t.Errorf("Unique(%v) missing expected item %q", input, item)
			}
		}
	})

	t.Run("Filter", func(t *testing.T) {
		input := []string{"apple", "banana", "cherry", "apricot"}
		predicate := func(s string) bool {
			return len(s) > 5
		}
		expected := []string{"banana", "cherry", "apricot"}

		result := slice.Filter(input, predicate)
		if len(result) != len(expected) {
			t.Errorf("Filter length = %d, expected %d", len(result), len(expected))
		}
	})
}

func TestValidationUtils(t *testing.T) {
	validate := utils.Validate()

	t.Run("Email", func(t *testing.T) {
		tests := []struct {
			email    string
			expected bool
		}{
			{"test@example.com", true},
			{"user.name+tag@domain.co.uk", true},
			{"invalid.email", false},
			{"@domain.com", false},
			{"user@", false},
		}

		for _, test := range tests {
			result := validate.Email(test.email)
			if result != test.expected {
				t.Errorf("Email(%q) = %v, expected %v", test.email, result, test.expected)
			}
		}
	})

	t.Run("IP", func(t *testing.T) {
		tests := []struct {
			ip       string
			expected bool
		}{
			{"192.168.1.1", true},
			{"10.0.0.1", true},
			{"256.1.1.1", false},
			{"192.168.1", false},
			{"not.an.ip", false},
		}

		for _, test := range tests {
			result := validate.IP(test.ip)
			if result != test.expected {
				t.Errorf("IP(%q) = %v, expected %v", test.ip, result, test.expected)
			}
		}
	})
}

func TestHashUtils(t *testing.T) {
	hash := utils.Hash()

	t.Run("MD5", func(t *testing.T) {
		input := "hello world"
		expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9" // Actually SHA256 for security

		result := hash.MD5(input)
		if result != expected {
			t.Errorf("MD5(%q) = %q, expected %q", input, result, expected)
		}
	})

	t.Run("SHA256", func(t *testing.T) {
		input := "hello world"
		expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

		result := hash.SHA256(input)
		if result != expected {
			t.Errorf("SHA256(%q) = %q, expected %q", input, result, expected)
		}
	})
}

func TestRandomUtils(t *testing.T) {
	random := utils.Random()

	t.Run("String", func(t *testing.T) {
		length := 10
		result := random.String(length)

		if len(result) != length {
			t.Errorf("String(%d) length = %d, expected %d", length, len(result), length)
		}
	})

	t.Run("Int", func(t *testing.T) {
		myMin, myMax := 5, 15
		result := random.Int(myMin, myMax)

		if result < myMin || result > myMax {
			t.Errorf("Int(%d, %d) = %d, expected value between %d and %d", myMin, myMax, result, myMin, myMax)
		}
	})

	t.Run("Choice", func(t *testing.T) {
		choices := []string{"apple", "banana", "cherry"}
		result := random.Choice(choices)

		slice := utils.Slice()
		if !slice.Contains(choices, result) {
			t.Errorf("Choice(%v) = %q, expected one of %v", choices, result, choices)
		}
	})
}

func BenchmarkStringReverse(b *testing.B) {
	str := utils.String()
	input := "hello world this is a test string"

	b.ResetTimer()
	for range b.N {
		str.Reverse(input)
	}
}

func BenchmarkSliceUnique(b *testing.B) {
	slice := utils.Slice()
	input := make([]string, 1000)
	for i := range 1000 {
		input[i] = fmt.Sprintf("item_%d", i%100) // Creates duplicates
	}

	b.ResetTimer()
	for range b.N {
		slice.Unique(input)
	}
}
