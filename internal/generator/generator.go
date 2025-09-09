// Package generator provides tool generation functionality for the toolbox project.
package generator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ToolType represents the type of tool to generate
type ToolType int

const (
	CLI ToolType = iota
	TUI
)

func (t ToolType) String() string {
	switch t {
	case CLI:
		return "CLI"
	case TUI:
		return "TUI"
	default:
		return "Unknown"
	}
}

// GeneratorModel represents the tool generator state
//
//nolint:revive // Using GeneratorModel instead of Model to avoid confusion with other model types
type GeneratorModel struct {
	step        int
	toolType    ToolType
	toolName    string
	toolDesc    string
	choices     []string
	cursor      int
	quitting    bool
	error       string
	success     string
	inputMode   bool
	inputText   strings.Builder
	inputPrompt string
}

// Generator styling
var (
	generatorTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 1)

	generatorItemStyle = lipgloss.NewStyle().
				PaddingLeft(4)

	generatorSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170"))

	generatorHelpStyle = lipgloss.NewStyle().
				PaddingLeft(4).
				PaddingTop(1).
				Foreground(lipgloss.Color("241"))

	generatorErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Bold(true)

	generatorSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("46")).
				Bold(true)

	generatorInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("33")).
				Background(lipgloss.Color("240")).
				Padding(0, 1)
)

// NewGeneratorModel creates a new generator model
func NewGeneratorModel() *GeneratorModel {
	return &GeneratorModel{
		step:    0,
		choices: []string{"CLI Tool", "TUI Tool", "Back to Main Menu"},
		cursor:  0,
	}
}

// Init implements tea.Model
func (m *GeneratorModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *GeneratorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.inputMode {
			return m.handleInputMode(msg)
		}
		return m.handleMenuMode(msg)
	}
	return m, nil
}

// handleInputMode handles input for tool name and description
func (m *GeneratorModel) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "esc":
		m.inputMode = false
		m.inputText.Reset()
		m.step--
		return m, nil
	case "enter":
		input := strings.TrimSpace(m.inputText.String())
		if input == "" {
			m.error = "Input cannot be empty"
			return m, nil
		}

		switch m.step {
		case 1: // Tool name
			m.toolName = input
			m.step++
			m.inputMode = true
			m.inputText.Reset()
			m.inputPrompt = "Enter tool description:"
			m.error = ""
		case 2: // Tool description
			m.toolDesc = input
			m.inputMode = false
			m.error = ""
			// Generate the tool
			if err := m.generateTool(); err != nil {
				m.error = fmt.Sprintf("Error generating tool: %v", err)
			} else {
				m.success = fmt.Sprintf("Successfully generated %s tool: %s", m.toolType.String(), m.toolName)
			}
			m.step++
		}
		return m, nil
	case "backspace":
		if m.inputText.Len() > 0 {
			content := m.inputText.String()
			m.inputText.Reset()
			m.inputText.WriteString(content[:len(content)-1])
		}
		return m, nil
	default:
		if len(msg.String()) == 1 {
			m.inputText.WriteString(msg.String())
		}
	}
	return m, nil
}

// handleMenuMode handles menu navigation
func (m *GeneratorModel) handleMenuMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
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
		return m.handleSelection()
	case "esc", "b":
		if m.step == 0 {
			// Return to main menu - this would need to be handled by parent
			return m, tea.Quit
		}
		// Go back to previous step
		m.step = 0
		m.cursor = 0
		m.error = ""
		m.success = ""
		m.inputMode = false
		m.inputText.Reset()
	case "r":
		if m.step == 3 {
			// Reset to create another tool
			m.step = 0
			m.cursor = 0
			m.error = ""
			m.success = ""
			m.toolName = ""
			m.toolDesc = ""
		}
	}
	return m, nil
}

// handleSelection handles menu item selection
func (m *GeneratorModel) handleSelection() (tea.Model, tea.Cmd) {
	switch m.step {
	case 0: // Tool type selection
		switch m.cursor {
		case 0:
			m.toolType = CLI
		case 1:
			m.toolType = TUI
		case 2:
			return m, tea.Quit // Back to main menu
		}
		m.step++
		m.inputMode = true
		m.inputPrompt = "Enter tool name (lowercase, no spaces):"
		m.error = ""
	}
	return m, nil
}

// View implements tea.Model
func (m *GeneratorModel) View() string {
	if m.quitting {
		return generatorTitleStyle.Render("Tool Generator") + "\n\nExiting...\n"
	}

	s := generatorTitleStyle.Render("🛠️  Go Tool Generator") + "\n\n"

	switch m.step {
	case 0: // Tool type selection
		s += "Select the type of tool to generate:\n\n"
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			if m.cursor == i {
				s += generatorSelectedStyle.Render(fmt.Sprintf("%s %s", cursor, choice))
			} else {
				s += generatorItemStyle.Render(fmt.Sprintf("%s %s", cursor, choice))
			}
			s += "\n"
		}
		s += generatorHelpStyle.Render("\nUse ↑/↓ or j/k to navigate, Enter to select, Esc to go back")

	case 1, 2: // Input mode
		s += fmt.Sprintf("Creating %s Tool\n\n", m.toolType.String())
		s += generatorItemStyle.Render(m.inputPrompt) + "\n"
		s += generatorInputStyle.Render(m.inputText.String()+"█") + "\n\n"

		if m.step == 1 {
			s += generatorHelpStyle.Render("Examples: filehasher, networkping, jsonformatter")
		} else {
			s += generatorHelpStyle.Render("Examples: A CLI tool for calculating file hashes")
		}
		s += "\n" + generatorHelpStyle.Render("Press Enter to continue, Esc to go back")

	case 3: // Completion
		s += "Tool Generation Complete!\n\n"
		if m.success != "" {
			s += generatorSuccessStyle.Render("✓ "+m.success) + "\n\n"
			s += generatorItemStyle.Render(fmt.Sprintf("Tool: %s", m.toolName)) + "\n"
			s += generatorItemStyle.Render(fmt.Sprintf("Type: %s", m.toolType.String())) + "\n"
			s += generatorItemStyle.Render(fmt.Sprintf("Description: %s", m.toolDesc)) + "\n\n"
			s += generatorItemStyle.Render("Files created:") + "\n"
			s += generatorItemStyle.Render(fmt.Sprintf("  • cmd/%s/%s/main.go", strings.ToLower(m.toolType.String()), m.toolName)) + "\n"
			s += generatorItemStyle.Render("  • README.md (updated)") + "\n"
			s += generatorItemStyle.Render("  • Makefile (updated)") + "\n\n"
			s += generatorHelpStyle.Render("Press 'r' to create another tool, 'b' to go back, or 'q' to quit")
		}
	}

	if m.error != "" {
		s += "\n" + generatorErrorStyle.Render("✗ "+m.error)
	}

	return s
}

// generateTool creates the actual tool files and directories
func (m *GeneratorModel) generateTool() error {
	// Validate tool name
	if !isValidToolName(m.toolName) {
		return errors.New("invalid tool name: use lowercase letters, numbers, and hyphens only")
	}

	// Create directory structure
	toolDir := filepath.Join("cmd", strings.ToLower(m.toolType.String()), m.toolName)
	if err := os.MkdirAll(toolDir, 0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate main.go file
	if err := m.generateMainFile(toolDir); err != nil {
		return fmt.Errorf("failed to generate main.go: %w", err)
	}

	// Generate additional files based on tool type
	switch m.toolType {
	case CLI:
		// CLI tools only need the main.go file, which is already generated
	case TUI:
		if err := m.generateTUIFiles(toolDir); err != nil {
			return fmt.Errorf("failed to generate TUI files: %w", err)
		}
	}

	// Update Makefile if needed
	if err := m.updateMakefile(); err != nil {
		return fmt.Errorf("failed to update Makefile: %w", err)
	}

	return nil
}

// generateMainFile creates the main.go file based on tool type
func (m *GeneratorModel) generateMainFile(toolDir string) error {
	var tmpl string

	switch m.toolType {
	case CLI:
		tmpl = cliTemplate
	case TUI:
		tmpl = tuiTemplate
	}

	t, err := template.New("main").Parse(tmpl)
	if err != nil {
		return err
	}

	data := struct {
		ToolName    string
		ToolDesc    string
		PackageName string
	}{
		ToolName:    m.toolName,
		ToolDesc:    m.toolDesc,
		PackageName: strings.ReplaceAll(m.toolName, "-", ""),
	}

	// #nosec G304 - This creates files in a controlled directory structure
	file, err := os.Create(filepath.Join(toolDir, "main.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, data)
}

// generateTUIFiles creates additional files for TUI tools
func (m *GeneratorModel) generateTUIFiles(_ string) error {
	// For now, TUI tools only need the main.go file
	// Could add additional model files here in the future
	return nil
}

// updateMakefile adds the new tool to the Makefile if needed
func (m *GeneratorModel) updateMakefile() error {
	// The current Makefile automatically discovers tools, so no update needed
	return nil
}

// isValidToolName checks if the tool name is valid
func isValidToolName(name string) bool {
	if name == "" {
		return false
	}
	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
			return false
		}
	}
	return true
}
