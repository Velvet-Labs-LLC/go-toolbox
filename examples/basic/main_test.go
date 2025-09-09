package main

import (
	"strings"
	"sync"
	"testing"

	"github.com/nate3d/toolbox/pkg/utils"
)

func BenchmarkStringUtils(b *testing.B) {
	str := utils.String()
	small := "hello world"
	large := strings.Repeat("a", 10000)

	b.Run("Reverse_Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = str.Reverse(small)
		}
	})
	b.Run("Reverse_Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = str.Reverse(large)
		}
	})
	b.Run("CamelCase", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = str.ToCamelCase(small)
		}
	})
	b.Run("SnakeCase", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = str.ToSnakeCase(small)
		}
	})
}

func BenchmarkSliceUtils(b *testing.B) {
	slice := utils.Slice()
	data := make([]string, 0, 1000)
	for i := 0; i < 1000; i++ {
		data = append(data, strings.Repeat("x", i%10))
	}

	b.Run("Unique", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = slice.Unique(data)
		}
	})
	b.Run("Sort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = slice.Sort(data)
		}
	})
	b.Run("Filter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = slice.Filter(data, func(s string) bool { return len(s) > 5 })
		}
	})
}

func BenchmarkRandomUtils(b *testing.B) {
	rand := utils.Random()

	b.Run("RandomString_10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = rand.String(10)
		}
	})
	b.Run("RandomInt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = rand.Int(1, 100)
		}
	})
}

func BenchmarkHashUtils(b *testing.B) {
	hash := utils.Hash()
	text := strings.Repeat("data", 100)

	b.Run("MD5", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = hash.MD5(text)
		}
	})
	b.Run("SHA256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = hash.SHA256(text)
		}
	})
}

func BenchmarkValidationUtils(b *testing.B) {
	val := utils.Validate()
	tests := []string{"test@example.com", "invalid@", "192.168.0.1", "999.999.999.999"}

	b.Run("EmailValidation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = val.Email(tests[i%len(tests)])
		}
	})
	b.Run("IPValidation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = val.IP(tests[i%len(tests)])
		}
	})
}

func BenchmarkFileUtils(b *testing.B) {
	file := utils.File()
	// go.mod lives two directories up from this example package
	existsTarget := "../../go.mod"

	b.Run("Exists", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = file.Exists(existsTarget)
		}
	})
}

func BenchmarkConvertUtils(b *testing.B) {
	conv := utils.Convert()

	b.Run("StringToInt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = conv.StringToInt("12345")
		}
	})
	b.Run("StringToFloat", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = conv.StringToFloat("123.456")
		}
	})
	b.Run("StringToBool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = conv.StringToBool("true")
		}
	})
}

func BenchmarkConcurrencyExample(b *testing.B) {
	hash := utils.Hash()
	urls := []string{"http://a", "http://b", "http://c"}

	b.Run("ConcurrentHash", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			wg.Add(len(urls))
			for _, u := range urls {
				go func(u string) {
					_ = hash.MD5(u)
					wg.Done()
				}(u)
			}
			wg.Wait()
		}
	})
}

// A simple smoke test to ensure main runs without panic
// TestMainFunction ensures main() runs without panicking
func TestMainFunction(t *testing.T) {
	// Just call main; output is ignored
	main()
}
