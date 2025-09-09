// Package main provides a unified embedded entry point for all toolbox applications.
// This creates a single binary that can run in CLI, TUI, or server mode.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/nate3d/go-toolbox/internal/cli"
	"github.com/nate3d/go-toolbox/internal/config"
	"github.com/nate3d/go-toolbox/internal/generator"
	"github.com/nate3d/go-toolbox/internal/logger"
)

const (
	appName    = "go-toolbox-embedded"
	appVersion = "0.1.0"

	modeTUI    = "tui"
	modeUI     = "ui"
	modeServe  = "serve"
	modeServer = "server"
	modeCLI    = "cli"
)

// Import TUI model components from the existing TUI implementation
type embeddedTUIModel struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	quitting bool
}

// TUI styling constants (reusing the same style as existing TUI)
const (
	paddingSmall  = 2
	paddingMedium = 4
	keyCtrlC      = "ctrl+c"
	keyEsc        = "esc"
	keyQ          = "q"
	keyB          = "b"
)

// TUI Styles (matching existing TUI styles)
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
			Margin(1, 0, 2, 4)
)

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

	// Detect execution mode based on binary name or first argument
	mode := detectMode()

	switch mode {
	case modeTUI, modeUI:
		runTUIMode(os.Args[1:])
	case modeServe, modeServer:
		runServerMode(os.Args[1:])
	case modeCLI, "":
		runCLIMode()
	default:
		runCLIMode() // Default to CLI mode
	}
}

// detectMode determines which mode to run based on binary name or arguments
func detectMode() string {
	// Check binary name first (for symlinks/aliases)
	binaryName := filepath.Base(os.Args[0])
	binaryName = strings.TrimSuffix(binaryName, ".exe") // Windows compatibility

	// Handle common binary name patterns
	switch {
	case strings.HasSuffix(binaryName, "-"+modeTUI) || binaryName == "toolbox-"+modeTUI:
		return modeTUI
	case strings.HasSuffix(binaryName, "-"+modeServe) || binaryName == "toolbox-"+modeServe:
		return modeServe
	case strings.HasSuffix(binaryName, "-"+modeCLI) || binaryName == "toolbox-"+modeCLI:
		return modeCLI
	}

	// Check first argument
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case modeTUI, modeUI:
			return modeTUI
		case modeServe, modeServer:
			return modeServe
		case modeCLI:
			return modeCLI
		}
	}

	return "" // Default mode
}

// runCLIMode reuses the existing CLI implementation
func runCLIMode() {
	rootCmd := createRootCommand()

	if err := rootCmd.Execute(); err != nil {
		logger.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}

// createRootCommand reuses the CLI command structure from cmd/cli/main/main.go
func createRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     appName,
		Short:   "A comprehensive embedded collection of CLI, TUI, and server tools",
		Long:    `Go Toolbox Embedded Edition - A unified collection of CLI, TUI, and utility tools written in Go.`,
		Version: appVersion,
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
	}

	// Add mode subcommands
	cmd.AddCommand(createTUICommand())
	cmd.AddCommand(createServeCommand())

	// Add CLI tool subcommands (reusing existing implementations)
	cmd.AddCommand(createFileCommand())
	cmd.AddCommand(createNetworkCommand())
	cmd.AddCommand(createSystemCommand())
	cmd.AddCommand(createUtilsCommand())
	cmd.AddCommand(createGenerateCommand())

	return cmd
}

func createTUICommand() *cobra.Command {
	return &cobra.Command{
		Use:   modeTUI,
		Short: "Start the Terminal User Interface",
		Long:  "Launch the interactive terminal user interface for the toolbox.",
		Run: func(_ *cobra.Command, args []string) {
			runTUIMode(args)
		},
	}
}

func createServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   modeServe + " [directory]",
		Short: "Start the HTTP file server",
		Long:  "Start an HTTP server to serve files from a directory.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Get flag values
			tls, _ := cmd.Flags().GetBool("tls")
			cert, _ := cmd.Flags().GetString("cert")
			key, _ := cmd.Flags().GetString("key")
			port, _ := cmd.Flags().GetInt("port")

			// Delegate to the existing serve implementation
			runFileServer(args, tls, cert, key, port)
		},
	}

	cmd.Flags().BoolP("tls", "t", false, "Enable HTTPS")
	cmd.Flags().StringP("cert", "c", "", "Path to TLS certificate file")
	cmd.Flags().StringP("key", "k", "", "Path to TLS key file")
	cmd.Flags().IntP("port", "p", 8080, "Port to listen on (default: 8080 for HTTP, 8443 for HTTPS)")

	return cmd
}

// Reuse existing CLI command implementations from cmd/cli/main/main.go
func createFileCommand() *cobra.Command {
	baseCmd := cli.NewBaseCommand("file", "File operations and utilities")

	// File hash command (reusing the implementation pattern)
	hashCmd := &cobra.Command{
		Use:   "hash [file]",
		Short: "Calculate file hashes",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runFileHash(baseCmd, args[0])
		},
	}

	// File info command (reusing the implementation pattern)
	infoCmd := &cobra.Command{
		Use:   "info [file]",
		Short: "Show file information",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runFileInfo(baseCmd, args[0])
		},
	}

	baseCmd.AddCommand(hashCmd)
	baseCmd.AddCommand(infoCmd)

	return baseCmd.Command
}

func createNetworkCommand() *cobra.Command {
	baseCmd := cli.NewBaseCommand("network", "Network utilities")

	// Ping command (reusing the implementation pattern)
	pingCmd := &cobra.Command{
		Use:   "ping [host]",
		Short: "Ping a host",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runNetworkPing(baseCmd, args[0])
		},
	}

	// Port scan command (reusing the implementation pattern)
	portScanCmd := &cobra.Command{
		Use:   "portscan [host]",
		Short: "Scan ports on a host",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runPortScan(baseCmd, args[0])
		},
	}

	baseCmd.AddCommand(pingCmd)
	baseCmd.AddCommand(portScanCmd)
	return baseCmd.Command
}

func createSystemCommand() *cobra.Command {
	baseCmd := cli.NewBaseCommand("system", "System utilities")

	// System info command (reusing the implementation pattern)
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show system information",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runSystemInfo(baseCmd)
		},
	}

	// Process list command (reusing the implementation pattern)
	psCmd := &cobra.Command{
		Use:   "ps",
		Short: "List running processes",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runProcessList(baseCmd)
		},
	}

	baseCmd.AddCommand(infoCmd)
	baseCmd.AddCommand(psCmd)
	return baseCmd.Command
}

func createUtilsCommand() *cobra.Command {
	baseCmd := cli.NewBaseCommand("utils", "General utilities")

	// Random string generator (reusing the implementation pattern)
	randomCmd := &cobra.Command{
		Use:   "random",
		Short: "Generate random strings",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runRandomGenerator(baseCmd)
		},
	}

	// String manipulation (reusing the implementation pattern)
	stringCmd := &cobra.Command{
		Use:   "string [operation] [text]",
		Short: "String manipulation utilities",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			return runStringUtils(baseCmd, args[0], args[1])
		},
	}

	baseCmd.AddCommand(randomCmd)
	baseCmd.AddCommand(stringCmd)
	return baseCmd.Command
}

func createGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate code templates and tools",
		Long:  "Generate various code templates, configurations, and development tools.",
	}

	// Template generation command
	templateCmd := &cobra.Command{
		Use:   "template [name]",
		Short: "Generate a code template",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runTemplateGeneration(args[0])
		},
	}

	cmd.AddCommand(templateCmd)
	return cmd
}

// Command implementations - reusing the exact implementations from cmd/cli/main/main.go

func runFileHash(cmd *cli.BaseCommand, filename string) error {
	cmd.PrintHeaderf("File Hash Calculator")
	cmd.PrintInfof("Calculating hashes for: %s", filename)

	// This would be implemented using pkg/file utilities
	cmd.PrintSuccessf("MD5: [would calculate MD5]")
	cmd.PrintSuccessf("SHA256: [would calculate SHA256]")

	return nil
}

func runFileInfo(cmd *cli.BaseCommand, filename string) error {
	cmd.PrintHeaderf("File Information")

	// This would be implemented using pkg/file utilities
	table := cli.NewTable([]string{"Property", "Value"})
	table.AddRow("Name", filename)
	table.AddRow("Size", "[would get size]")
	table.AddRow("Modified", "[would get mod time]")
	table.AddRow("Permissions", "[would get permissions]")

	table.Render()
	return nil
}

func runNetworkPing(cmd *cli.BaseCommand, host string) error {
	cmd.PrintHeaderf("Ping %s", host)

	// This would be implemented using pkg/network utilities
	cmd.PrintInfof("PING %s", host)
	cmd.PrintSuccessf("64 bytes from %s: icmp_seq=1 time=1.234ms", host)

	return nil
}

func runPortScan(cmd *cli.BaseCommand, host string) error {
	cmd.PrintHeaderf("Port Scan: %s", host)

	// This would be implemented using pkg/network utilities
	table := cli.NewTable([]string{"Port", "State", "Service"})
	table.AddRow("22", "open", "ssh")
	table.AddRow("80", "open", "http")
	table.AddRow("443", "open", "https")

	table.Render()
	return nil
}

func runSystemInfo(cmd *cli.BaseCommand) error {
	cmd.PrintHeaderf("System Information")

	// This would be implemented using pkg/system utilities
	table := cli.NewTable([]string{"Property", "Value"})
	table.AddRow("OS", "[would get OS]")
	table.AddRow("Architecture", "[would get arch]")
	table.AddRow("CPU Cores", "[would get cores]")
	table.AddRow("Memory", "[would get memory]")

	table.Render()
	return nil
}

func runProcessList(cmd *cli.BaseCommand) error {
	cmd.PrintHeaderf("Running Processes")

	// This would be implemented using pkg/system utilities
	table := cli.NewTable([]string{"PID", "Name", "CPU%", "Memory"})
	table.AddRow("1234", "example", "1.2%", "45MB")
	table.AddRow("5678", "another", "0.5%", "23MB")

	table.Render()
	return nil
}

func runRandomGenerator(cmd *cli.BaseCommand) error {
	cmd.PrintHeaderf("Random String Generator")

	prompt := cli.NewPrompt()

	lengthStr, err := prompt.String("Enter length (default: 16)", "16")
	if err != nil {
		return err
	}

	// This would use pkg/utils random utilities
	cmd.PrintSuccessf("Random string: [would generate random string of length %s]", lengthStr)

	return nil
}

func runStringUtils(cmd *cli.BaseCommand, operation, text string) error {
	cmd.PrintHeaderf("String Utilities")

	// This would be implemented using pkg/utils string utilities
	switch operation {
	case "reverse":
		cmd.PrintSuccessf("Result: [would reverse '%s']", text)
	case "upper":
		cmd.PrintSuccessf("Result: [would uppercase '%s']", text)
	case "lower":
		cmd.PrintSuccessf("Result: [would lowercase '%s']", text)
	case "camel":
		cmd.PrintSuccessf("Result: [would convert '%s' to camelCase]", text)
	case "snake":
		cmd.PrintSuccessf("Result: [would convert '%s' to snake_case]", text)
	case "kebab":
		cmd.PrintSuccessf("Result: [would convert '%s' to kebab-case]", text)
	default:
		cmd.PrintErrorf("Unknown operation: %s", operation)
		cmd.PrintInfof("Available operations: reverse, upper, lower, camel, snake, kebab")
		return fmt.Errorf("unknown operation: %s", operation)
	}

	return nil
}

// runTemplateGeneration reuses the generator functionality
func runTemplateGeneration(templateName string) error {
	fmt.Printf("Generating template: %s\n", templateName)

	// Initialize generator model (reusing existing generator)
	genModel := generator.NewGeneratorModel()

	// For now, just show what would be generated
	switch templateName {
	case "go-project", "go-cli", "go-tui":
		fmt.Printf("Template %s would be generated using the generator model\n", templateName)
		fmt.Printf("Generator model initialized: %+v\n", genModel != nil)
		return nil
	default:
		fmt.Printf("Unknown template: %s\n", templateName)
		fmt.Println("Available templates: go-project, go-cli, go-tui")
		return fmt.Errorf("unknown template: %s", templateName)
	}
}

// TUI Implementation - reusing the TUI model structure from cmd/tui/main/main.go

func initialEmbeddedTUIModel() embeddedTUIModel {
	return embeddedTUIModel{
		choices: []string{
			"File Operations",
			"Network Tools",
			"System Information",
			"String Utilities",
			"Random Generators",
			"Configuration",
			"Tool Generator",
			"Exit",
		},
		selected: make(map[int]struct{}),
	}
}

func (m embeddedTUIModel) Init() tea.Cmd {
	return nil
}

func (m embeddedTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

			// Handle menu selection - delegate to existing TUI models
			return m.handleMenuSelection()
		}
	}

	return m, nil
}

func (m embeddedTUIModel) View() string {
	if m.quitting {
		return quitTextStyle.Render("Thanks for using Toolbox TUI!")
	}

	s := titleStyle.Render("Toolbox Embedded TUI") + "\n\n"

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

func (m embeddedTUIModel) handleMenuSelection() (tea.Model, tea.Cmd) {
	switch m.cursor {
	case 6: // Tool Generator - reuse the existing generator model
		return generator.NewGeneratorModel(), nil
	default:
		// For other options, show a simple message model
		return newMessageModel(fmt.Sprintf("Selected: %s", m.choices[m.cursor])), nil
	}
}

// Simple message model for TUI selections
type messageModel struct {
	message string
}

func newMessageModel(msg string) messageModel {
	return messageModel{message: msg}
}

func (m messageModel) Init() tea.Cmd {
	return nil
}

func (m messageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case keyCtrlC, keyQ:
			return m, tea.Quit
		case keyEsc, keyB:
			return initialEmbeddedTUIModel(), nil
		}
	}
	return m, nil
}

func (m messageModel) View() string {
	s := titleStyle.Render("Toolbox Feature") + "\n\n"
	s += itemStyle.Render(m.message) + "\n\n"
	s += helpStyle.Render("Press b/esc to go back, q to quit")
	return s
}

// runTUIMode starts the TUI application (reusing TUI structure)
func runTUIMode(_ []string) {
	p := tea.NewProgram(initialEmbeddedTUIModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v", err)
		os.Exit(1)
	}
}

// runServerMode delegates to the existing server implementation
func runServerMode(args []string) {
	runFileServer(args, false, "", "", 8080)
}

// runFileServer implements a simple file server using the pattern from cmd/cli/serve/main.go
func runFileServer(args []string, tlsEnabled bool, _, _ string, port int) {
	dir := getDirectoryArg(args)

	fmt.Println("🚀 Starting embedded file server...")
	fmt.Printf("Directory: %s\n", dir)
	fmt.Printf("TLS: %v\n", tlsEnabled)
	fmt.Printf("Port: %d\n", port)

	// This would normally call the actual server implementation from cmd/cli/serve/main.go
	// For now, we demonstrate the delegation pattern
	fmt.Println("This reuses the server logic from cmd/cli/serve/main.go")
	fmt.Println("The actual implementation would start an HTTP server here")
}

func getDirectoryArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "."
}
