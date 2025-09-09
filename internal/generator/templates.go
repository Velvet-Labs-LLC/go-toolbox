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
