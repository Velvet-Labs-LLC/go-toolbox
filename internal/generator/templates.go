package generator

// CLI Tool Template
const cliTemplate = `// Package main provides {{.ToolDesc}}
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nate3d/go-toolbox/internal/config"
	"github.com/nate3d/go-toolbox/internal/logger"
)

const appName = "{{.ToolName}}"

func main() {
	// Command line flags
	var (
		version = flag.Bool("version", false, "Show version information")
		help    = flag.Bool("help", false, "Show help information")
		verbose = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *version {
		fmt.Printf("%s version 1.0.0\n", appName)
		return
	}

	// Initialize configuration
	if err := config.Init(appName); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logLevel := "info"
	if *verbose {
		logLevel = "debug"
	}

	logConfig := logger.Config{
		Level:      logger.LogLevel(logLevel),
		Output:     config.GetString("log_file"),
		Format:     "text",
		WithCaller: false,
		WithTime:   true,
	}
	if err := logger.Init(logConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	logger.Info("Starting {{.ToolName}}")

	// TODO: Implement your CLI tool logic here
	fmt.Printf("{{.ToolDesc}}\n")
	fmt.Printf("This is a generated CLI tool template.\n")
	fmt.Printf("Implement your functionality in the main() function.\n")

	logger.Info("{{.ToolName}} completed successfully")
}

func showHelp() {
	fmt.Printf("{{.ToolName}} - {{.ToolDesc}}\n\n")
	fmt.Printf("Usage: %s [options]\n\n", appName)
	fmt.Printf("Options:\n")
	flag.PrintDefaults()
	fmt.Printf("\nGenerated with Go Toolbox Generator\n")
}
`

// TUI Tool Template
const tuiTemplate = `// Package main provides {{.ToolDesc}}
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/nate3d/go-toolbox/internal/config"
	"github.com/nate3d/go-toolbox/internal/logger"
)

const appName = "{{.ToolName}}"

// Styles
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
		PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
		PaddingLeft(2).
		Foreground(lipgloss.Color("170"))

	helpStyle = lipgloss.NewStyle().
		PaddingLeft(4).
		PaddingTop(1).
		Foreground(lipgloss.Color("241"))

	quitTextStyle = lipgloss.NewStyle().
		Margin(1, 0, 2, 4)
)

// Model represents the application state
type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	quitting bool
}

// initialModel creates the initial model
func initialModel() model {
	return model{
		choices: []string{
			"Option 1",
			"Option 2",
			"Option 3",
			"Exit",
		},
		selected: make(map[int]struct{}),
	}
}

// Init implements tea.Model
func (m model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			if m.cursor == len(m.choices)-1 { // Exit option
				m.quitting = true
				return m, tea.Quit
			}

			// TODO: Handle menu selection here
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

// View implements tea.Model
func (m model) View() string {
	if m.quitting {
		return quitTextStyle.Render("Thanks for using {{.ToolName}}!")
	}

	s := titleStyle.Render("{{.ToolName}}") + "\n\n"
	s += itemStyle.Render("{{.ToolDesc}}") + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		if m.cursor == i {
			s += selectedItemStyle.Render(fmt.Sprintf("%s [%s] %s", cursor, checked, choice))
		} else {
			s += itemStyle.Render(fmt.Sprintf("%s [%s] %s", cursor, checked, choice))
		}
		s += "\n"
	}

	s += helpStyle.Render("\nPress q to quit, ↑/↓ or j/k to navigate, enter to select.")

	return s
}

func main() {
	// Initialize configuration
	if err := config.Init(appName); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logConfig := logger.Config{
		Level:      logger.LogLevel(config.GetString("log_level")),
		Output:     config.GetString("log_file"),
		Format:     "text",
		WithCaller: false,
		WithTime:   true,
	}
	if err := logger.Init(logConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	logger.Info("Starting {{.ToolName}} TUI")

	// Start the TUI
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Error("TUI application failed", "error", err)
		os.Exit(1)
	}
}
`

// Web Tool Template
const webTemplate = `// Package main provides {{.ToolDesc}}
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nate3d/go-toolbox/internal/config"
	"github.com/nate3d/go-toolbox/internal/logger"
)

const appName = "{{.ToolName}}"

type server struct {
	templates *template.Template
}

func main() {
	// Initialize configuration
	if err := config.Init(appName); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logConfig := logger.Config{
		Level:      logger.LogLevel(config.GetString("log_level")),
		Output:     config.GetString("log_file"),
		Format:     "text",
		WithCaller: false,
		WithTime:   true,
	}
	if err := logger.Init(logConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	// Initialize server
	srv := &server{}

	// Load templates
	templatePath := filepath.Join("cmd", "web", "{{.ToolName}}", "templates", "*.html")
	var err error
	srv.templates, err = template.ParseGlob(templatePath)
	if err != nil {
		logger.Error("Failed to load templates", "error", err, "path", templatePath)
		// Create a simple inline template as fallback
		srv.templates = template.Must(template.New("index").Parse(defaultIndexTemplate))
	}

	// Setup routes
	http.HandleFunc("/", srv.handleIndex)
	http.HandleFunc("/api/status", srv.handleAPIStatus)

	// Serve static files
	staticPath := filepath.Join("cmd", "web", "{{.ToolName}}", "static")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	port := config.GetString("port")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting {{.ToolName}} web server", "port", port)
	fmt.Printf("{{.ToolDesc}}\n")
	fmt.Printf("Server starting on http://localhost:%s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title       string
		Description string
		AppName     string
	}{
		Title:       "{{.ToolName}}",
		Description: "{{.ToolDesc}}",
		AppName:     appName,
	}

	if err := s.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		logger.Error("Failed to execute template", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *server) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, ` + "`" + `{"status": "ok", "service": "%s", "description": "%s"}` + "`" + `, appName, "{{.ToolDesc}}")
}

const defaultIndexTemplate = ` + "`" + `<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 800px; margin: 0 auto; }
        .header { background: #7D56F4; color: white; padding: 20px; border-radius: 8px; }
        .content { margin: 20px 0; }
        .footer { margin-top: 40px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.Title}}</h1>
            <p>{{.Description}}</p>
        </div>
        <div class="content">
            <h2>Welcome to {{.AppName}}</h2>
            <p>This is a generated web application template.</p>
            <p>Implement your web functionality by modifying the handlers and templates.</p>

            <h3>API Endpoints:</h3>
            <ul>
                <li><a href="/api/status">GET /api/status</a> - Service status</li>
            </ul>

            <h3>Quick Start:</h3>
            <ol>
                <li>Modify the templates in <code>templates/</code></li>
                <li>Add static assets to <code>static/</code></li>
                <li>Implement your handlers in <code>main.go</code></li>
                <li>Add new routes as needed</li>
            </ol>
        </div>
        <div class="footer">
            Generated with Go Toolbox Generator
        </div>
    </div>
</body>
</html>` + "`" + `
`

// HTML Template for web tools
const htmlTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 40px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header {
            background: #7D56F4;
            color: white;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 30px;
        }
        .content {
            line-height: 1.6;
        }
        .footer {
            margin-top: 40px;
            color: #666;
            font-size: 12px;
            text-align: center;
            border-top: 1px solid #eee;
            padding-top: 20px;
        }
        code {
            background: #f0f0f0;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
        }
        .api-list {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 5px;
            border-left: 4px solid #7D56F4;
        }
        .quick-start {
            background: #e8f5e8;
            padding: 15px;
            border-radius: 5px;
            border-left: 4px solid #28a745;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.Title}}</h1>
            <p>{{.Description}}</p>
        </div>
        <div class="content">
            <h2>🎉 Welcome to {{.AppName}}</h2>
            <p>This is a generated web application template created with the Go Toolbox Generator.</p>

            <div class="api-list">
                <h3>📡 API Endpoints:</h3>
                <ul>
                    <li><a href="/api/status">GET /api/status</a> - Service status check</li>
                    <li><a href="/static/">GET /static/</a> - Static file serving</li>
                </ul>
            </div>

            <div class="quick-start">
                <h3>🚀 Quick Start Guide:</h3>
                <ol>
                    <li>Modify templates in <code>templates/</code> directory</li>
                    <li>Add CSS, JS, and images to <code>static/</code> directory</li>
                    <li>Implement your business logic in <code>main.go</code></li>
                    <li>Add new HTTP routes and handlers as needed</li>
                    <li>Update configuration in <code>configs/config.yaml</code></li>
                </ol>
            </div>

            <h3>🛠️ Development Tips:</h3>
            <ul>
                <li>Use the logger for structured logging</li>
                <li>Configuration is loaded from the config package</li>
                <li>Templates are automatically reloaded in development</li>
                <li>Static files are served from the /static/ route</li>
            </ul>
        </div>
        <div class="footer">
            Generated with ❤️ by Go Toolbox Generator
        </div>
    </div>
</body>
</html>`
