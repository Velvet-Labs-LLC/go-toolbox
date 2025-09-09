// Package main provides a sample TUI application using Bubble Tea.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/nate3d/go-toolbox/internal/config"
	"github.com/nate3d/go-toolbox/internal/generator"
	"github.com/nate3d/go-toolbox/internal/logger"
)

const appName = "toolbox-tui"

// Constants for UI styling
const (
	// Padding values
	paddingSmall  = 2
	paddingMedium = 4

	// Key bindings
	keyCtrlC = "ctrl+c"
	keyEsc   = "esc"
	keyQ     = "q"
	keyB     = "b"
)

// Styles.
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(paddingMedium)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(paddingSmall).
				Foreground(lipgloss.Color("170"))

	helpStyle = lipgloss.NewStyle().
			PaddingLeft(paddingMedium).
			PaddingTop(1).
			Foreground(lipgloss.Color("241"))

	quitTextStyle = lipgloss.NewStyle().
			Margin(1, 0, paddingSmall, paddingMedium)
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
			"File Operations",
			"Network Tools",
			"System Information",
			"String Utilities",
			"Random Generators",
			"Configuration",
			"🛠️  Tool Generator",
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
		case keyCtrlC, keyQ:
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

			// Handle menu selection
			return m.handleMenuSelection()
		}
	}

	return m, nil
}

// View implements tea.Model
func (m model) View() string {
	if m.quitting {
		return quitTextStyle.Render("Thanks for using Toolbox TUI!")
	}

	s := titleStyle.Render("Toolbox TUI") + "\n\n"

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

	s += helpStyle.Render("\nPress q to quit, enter to select.")

	return s
}

// handleMenuSelection handles menu item selection
func (m model) handleMenuSelection() (tea.Model, tea.Cmd) {
	switch m.cursor {
	case 0: // File Operations
		return NewFileOperationsModel(), nil
	case 1: // Network Tools
		return NewNetworkToolsModel(), nil
	case 2: // System Information
		return NewSystemInfoModel(), nil
	case 3: // String Utilities
		return NewStringUtilsModel(), nil
	case 4: // Random Generators
		return NewRandomGenModel(), nil
	case 5: // Configuration
		return NewConfigModel(), nil
	case 6: // Tool Generator
		return generator.NewGeneratorModel(), nil
	}
	return m, nil
}

// File Operations Model
type fileOpsModel struct{}

func NewFileOperationsModel() tea.Model {
	return fileOpsModel{}
}

func (m fileOpsModel) Init() tea.Cmd {
	return nil
}

func (m fileOpsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case keyCtrlC, "q":
			return m, tea.Quit
		case keyEsc, keyB:
			return initialModel(), nil
		}
	}
	return m, nil
}

func (m fileOpsModel) View() string {
	s := titleStyle.Render("File Operations") + "\n\n"
	s += itemStyle.Render("This is where file operations would be implemented.") + "\n"
	s += itemStyle.Render("Features could include:") + "\n"
	s += itemStyle.Render("  • File hash calculation") + "\n"
	s += itemStyle.Render("  • File size analysis") + "\n"
	s += itemStyle.Render("  • Directory tree view") + "\n"
	s += itemStyle.Render("  • File search") + "\n\n"
	s += helpStyle.Render("Press 'b' or 'esc' to go back, 'q' to quit.")
	return s
}

// Network Tools Model
type networkToolsModel struct{}

func NewNetworkToolsModel() tea.Model {
	return networkToolsModel{}
}

func (m networkToolsModel) Init() tea.Cmd {
	return nil
}

func (m networkToolsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case keyCtrlC, "q":
			return m, tea.Quit
		case keyEsc, keyB:
			return initialModel(), nil
		}
	}
	return m, nil
}

func (m networkToolsModel) View() string {
	s := titleStyle.Render("Network Tools") + "\n\n"
	s += itemStyle.Render("Network utilities would be implemented here.") + "\n"
	s += itemStyle.Render("Features could include:") + "\n"
	s += itemStyle.Render("  • Ping tool") + "\n"
	s += itemStyle.Render("  • Port scanner") + "\n"
	s += itemStyle.Render("  • Network interface info") + "\n"
	s += itemStyle.Render("  • DNS lookup") + "\n\n"
	s += helpStyle.Render("Press 'b' or 'esc' to go back, 'q' to quit.")
	return s
}

// System Information Model
type systemInfoModel struct{}

func NewSystemInfoModel() tea.Model {
	return systemInfoModel{}
}

func (m systemInfoModel) Init() tea.Cmd {
	return nil
}

func (m systemInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case keyCtrlC, "q":
			return m, tea.Quit
		case keyEsc, keyB:
			return initialModel(), nil
		}
	}
	return m, nil
}

func (m systemInfoModel) View() string {
	s := titleStyle.Render("System Information") + "\n\n"
	s += itemStyle.Render("System information would be displayed here.") + "\n"
	s += itemStyle.Render("Information could include:") + "\n"
	s += itemStyle.Render("  • OS and version") + "\n"
	s += itemStyle.Render("  • CPU information") + "\n"
	s += itemStyle.Render("  • Memory usage") + "\n"
	s += itemStyle.Render("  • Disk usage") + "\n"
	s += itemStyle.Render("  • Running processes") + "\n\n"
	s += helpStyle.Render("Press 'b' or 'esc' to go back, 'q' to quit.")
	return s
}

// String Utilities Model
type stringUtilsModel struct{}

func NewStringUtilsModel() tea.Model {
	return stringUtilsModel{}
}

func (m stringUtilsModel) Init() tea.Cmd {
	return nil
}

func (m stringUtilsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case keyCtrlC, "q":
			return m, tea.Quit
		case keyEsc, keyB:
			return initialModel(), nil
		}
	}
	return m, nil
}

func (m stringUtilsModel) View() string {
	s := titleStyle.Render("String Utilities") + "\n\n"
	s += itemStyle.Render("String manipulation tools would be here.") + "\n"
	s += itemStyle.Render("Operations could include:") + "\n"
	s += itemStyle.Render("  • Case conversions") + "\n"
	s += itemStyle.Render("  • String reversal") + "\n"
	s += itemStyle.Render("  • Text encoding/decoding") + "\n"
	s += itemStyle.Render("  • Regular expression testing") + "\n\n"
	s += helpStyle.Render("Press 'b' or 'esc' to go back, 'q' to quit.")
	return s
}

// Random Generator Model
type randomGenModel struct{}

func NewRandomGenModel() tea.Model {
	return randomGenModel{}
}

func (m randomGenModel) Init() tea.Cmd {
	return nil
}

func (m randomGenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case keyCtrlC, "q":
			return m, tea.Quit
		case keyEsc, keyB:
			return initialModel(), nil
		}
	}
	return m, nil
}

func (m randomGenModel) View() string {
	s := titleStyle.Render("Random Generators") + "\n\n"
	s += itemStyle.Render("Random generation tools would be here.") + "\n"
	s += itemStyle.Render("Generators could include:") + "\n"
	s += itemStyle.Render("  • Random strings") + "\n"
	s += itemStyle.Render("  • UUIDs") + "\n"
	s += itemStyle.Render("  • Passwords") + "\n"
	s += itemStyle.Render("  • Random numbers") + "\n\n"
	s += helpStyle.Render("Press 'b' or 'esc' to go back, 'q' to quit.")
	return s
}

// Configuration Model
type configModel struct{}

func NewConfigModel() tea.Model {
	return configModel{}
}

func (m configModel) Init() tea.Cmd {
	return nil
}

func (m configModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case keyCtrlC, "q":
			return m, tea.Quit
		case keyEsc, keyB:
			return initialModel(), nil
		}
	}
	return m, nil
}

func (m configModel) View() string {
	s := titleStyle.Render("Configuration") + "\n\n"
	s += itemStyle.Render("Configuration settings would be here.") + "\n"
	s += itemStyle.Render("Settings could include:") + "\n"
	s += itemStyle.Render("  • Theme selection") + "\n"
	s += itemStyle.Render("  • Default output formats") + "\n"
	s += itemStyle.Render("  • Logging preferences") + "\n"
	s += itemStyle.Render("  • Key bindings") + "\n\n"
	s += helpStyle.Render("Press 'b' or 'esc' to go back, 'q' to quit.")
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

	// Start the TUI
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Error("TUI application failed", "error", err)
		os.Exit(1)
	}
}
