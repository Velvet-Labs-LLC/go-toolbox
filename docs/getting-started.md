# Getting Started with Go Toolbox

This guide will help you get started with the Go Toolbox project, especially if you're coming from a Python background.

## Go vs Python: Key Differences

### Language Characteristics

| Aspect | Go | Python |
|--------|----|---------|
| **Type System** | Statically typed | Dynamically typed |
| **Compilation** | Compiled to native binary | Interpreted |
| **Memory Management** | Garbage collected | Garbage collected |
| **Concurrency** | Built-in goroutines & channels | Threading/asyncio |
| **Performance** | Very fast | Moderate |
| **Deployment** | Single binary | Requires interpreter |

### Syntax Comparison

#### Variables and Types

**Python:**

```python
name = "John"
age = 30
is_active = True
numbers = [1, 2, 3, 4, 5]
```

**Go:**

```go
var name string = "John"
// or
name := "John"

var age int = 30
// or  
age := 30

var isActive bool = true
// or
isActive := true

var numbers []int = []int{1, 2, 3, 4, 5}
// or
numbers := []int{1, 2, 3, 4, 5}
```

#### Functions

**Python:**

```python
def greet(name: str) -> str:
    return f"Hello, {name}!"

def add_numbers(a: int, b: int) -> int:
    return a + b
```

**Go:**

```go
func greet(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}

func addNumbers(a, b int) int {
    return a + b
}

// Multiple returns (unique to Go)
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

#### Error Handling

**Python:**

```python
try:
    result = risky_operation()
except Exception as e:
    print(f"Error: {e}")
    return None
return result
```

**Go:**

```go
result, err := riskyOperation()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
// Use result
```

#### Classes vs Structs

**Python:**

```python
class Person:
    def __init__(self, name: str, age: int):
        self.name = name
        self.age = age
    
    def greet(self) -> str:
        return f"Hi, I'm {self.name}"

person = Person("Alice", 25)
print(person.greet())
```

**Go:**

```go
type Person struct {
    Name string
    Age  int
}

func (p Person) Greet() string {
    return fmt.Sprintf("Hi, I'm %s", p.Name)
}

person := Person{Name: "Alice", Age: 25}
fmt.Println(person.Greet())
```

### Package Management

**Python:**

```bash
# Create virtual environment
python -m venv venv
source venv/bin/activate  # or venv\Scripts\activate on Windows

# Install packages
pip install requests numpy

# Requirements file
pip freeze > requirements.txt
pip install -r requirements.txt
```

**Go:**

```bash
# Initialize module
go mod init github.com/username/project

# Add dependencies (automatic)
go get github.com/spf13/cobra

# Update dependencies
go mod tidy

# Download dependencies
go mod download
```

## Go Toolbox Project Structure

Our project follows Go's standard project layout:

```(text)
toolbox/
├── cmd/                    # Main applications
│   ├── cli/               # CLI tools
│   ├── tui/               # Terminal UI apps
│   └── web/               # Web applications
├── internal/              # Private code (like Python's private modules)
│   ├── config/            # Configuration management
│   ├── logger/            # Logging utilities
│   └── common/            # Shared internal code
├── pkg/                   # Public libraries (like Python packages)
│   ├── utils/             # General utilities
│   ├── network/           # Network utilities
│   └── file/              # File utilities
├── configs/               # Configuration files
├── scripts/               # Build scripts
├── test/                  # Test utilities
├── docs/                  # Documentation
└── examples/              # Usage examples
```

## Building and Running

### Development Commands

```bash
# Run without building
go run ./cmd/cli/main

# Build single binary
go build -o bin/toolbox ./cmd/cli/main

# Build with optimizations
go build -ldflags="-s -w" -o bin/toolbox ./cmd/cli/main

# Cross-compile for different platforms
GOOS=windows GOARCH=amd64 go build -o bin/toolbox.exe ./cmd/cli/main
GOOS=linux GOARCH=amd64 go build -o bin/toolbox-linux ./cmd/cli/main
```

### Using Make (like Python's Makefile)

```bash
# Build all applications
make build

# Run tests
make test

# Run linters
make lint

# Clean build artifacts
make clean

# Install development tools
make dev-setup
```

## Testing

Go has built-in testing, similar to Python's unittest:

**Python:**

```python
import unittest

class TestUtils(unittest.TestCase):
    def test_add(self):
        self.assertEqual(add(2, 3), 5)
    
    def test_divide(self):
        with self.assertRaises(ZeroDivisionError):
            divide(5, 0)
```

**Go:**

```go
package main

import "testing"

func TestAdd(t *testing.T) {
    result := add(2, 3)
    if result != 5 {
        t.Errorf("add(2, 3) = %d, expected 5", result)
    }
}

func TestDivide(t *testing.T) {
    _, err := divide(5, 0)
    if err == nil {
        t.Error("divide(5, 0) should return an error")
    }
}
```

## Concurrency

Go's concurrency model is different from Python:

**Python (asyncio):**

```python
import asyncio

async def fetch_url(url):
    # async HTTP request
    return await http_get(url)

async def main():
    tasks = [fetch_url(url) for url in urls]
    results = await asyncio.gather(*tasks)
```

**Go (goroutines):**

```go
func fetchURL(url string) string {
    // HTTP request
    return httpGet(url)
}

func main() {
    var wg sync.WaitGroup
    results := make(chan string, len(urls))
    
    for _, url := range urls {
        wg.Add(1)
        go func(u string) {
            defer wg.Done()
            results <- fetchURL(u)
        }(url)
    }
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    for result := range results {
        fmt.Println(result)
    }
}
```

## Configuration Management

Our project uses Viper for configuration, similar to Python's configparser or pydantic:

**Python:**

```python
import yaml

with open('config.yaml') as f:
    config = yaml.safe_load(f)

database_url = config['database']['url']
```

**Go:**

```go
import "github.com/spf13/viper"

viper.SetConfigName("config")
viper.AddConfigPath("./configs")
viper.ReadInConfig()

databaseURL := viper.GetString("database.url")
```

## Logging

We use structured logging similar to Python's logging module:

**Python:**

```python
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

logger.info("Processing file", extra={"filename": "data.txt"})
logger.error("Failed to process", exc_info=True)
```

**Go:**

```go
import "log/slog"

logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

logger.Info("Processing file", "filename", "data.txt")
logger.Error("Failed to process", "error", err)
```

## Next Steps

1. **Install Go**: Download from <https://golang.org/dl/>
2. **Set up your editor**: VS Code with Go extension is recommended
3. **Clone and build**:

   ```bash
   git clone <repo>
   cd toolbox
   go mod download
   make build
   ```

4. **Explore the code**: Start with `cmd/cli/main/main.go`
5. **Add your first tool**: Create a new command in the CLI application

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go by Example](https://gobyexample.com/)
- [Go Modules Reference](https://golang.org/ref/mod)
- [Standard Library](https://pkg.go.dev/std)
