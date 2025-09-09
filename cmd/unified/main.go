// Package main provides a unified entry point for all toolbox applications.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	appName    = "go-toolbox"
	appVersion = "0.1.0"

	modeTUI    = "tui"
	modeUI     = "ui"
	modeServe  = "serve"
	modeServer = "server"
	modeCLI    = "cli"
)

func main() {
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

// runCLIMode starts the CLI interface with subcommands
func runCLIMode() {
	rootCmd := &cobra.Command{
		Use:     appName,
		Short:   "Go Toolbox - A collection of useful Go tools",
		Version: appVersion,
		Long: `Go Toolbox is a unified collection of CLI and TUI tools for Go development.

Available modes:
  cli     - Command-line interface (default)
  tui     - Terminal user interface  
  serve   - HTTP server mode

Examples:
  go-toolbox                    # CLI mode (default)
  go-toolbox tui               # TUI mode
  go-toolbox serve             # Server mode
  go-toolbox generate --help   # CLI subcommand
  
You can also create symlinks for convenience:
  ln -s go-toolbox toolbox-tui
  ln -s go-toolbox toolbox-serve`,
	}

	// Add mode subcommands
	rootCmd.AddCommand(createTUICommand())
	rootCmd.AddCommand(createServeCommand())
	rootCmd.AddCommand(createGenerateCommand())
	rootCmd.AddCommand(createVersionCommand())

	// Add global flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file path")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().Bool("debug", false, "debug mode")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// createTUICommand creates the TUI subcommand
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

// createServeCommand creates the serve subcommand
func createServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   modeServe,
		Short: "Start the HTTP server",
		Long:  "Start the HTTP server component of the toolbox.",
		Run: func(_ *cobra.Command, args []string) {
			runServerMode(args)
		},
	}

	// Add server-specific flags
	cmd.Flags().StringP("addr", "a", ":8080", "server address")
	cmd.Flags().StringP("cert", "", "", "TLS certificate file")
	cmd.Flags().StringP("key", "", "", "TLS key file")
	cmd.Flags().BoolP("tls", "t", false, "enable TLS with auto-generated certificates")

	return cmd
}

// createGenerateCommand creates CLI commands for code generation
func createGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate code templates and tools",
		Long:  "Generate various code templates, configurations, and development tools.",
	}

	// Add generate subcommands
	cmd.AddCommand(&cobra.Command{
		Use:   "template [name]",
		Short: "Generate a code template",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			fmt.Printf("Generating template: %s\n", args[0])
			// TODO: Implement template generation
		},
	})

	return cmd
}

// createVersionCommand creates the version command
func createVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("%s version %s\n", appName, appVersion)
		},
	}
}

// runTUIMode starts the TUI application
func runTUIMode(args []string) {
	fmt.Println("🎨 Starting TUI mode...")
	fmt.Println("This would launch the Terminal User Interface")
	fmt.Printf("Arguments: %v\n", args)

	// TODO: Import and call the actual TUI code
	// This is where you'd call the code from cmd/tui/main/main.go
	fmt.Println("\nTo implement: Move TUI logic from cmd/tui/main/main.go here")
}

// runServerMode starts the server application
func runServerMode(args []string) {
	fmt.Println("🚀 Starting server mode...")
	fmt.Println("This would start the HTTP server")
	fmt.Printf("Arguments: %v\n", args)

	// TODO: Import and call the actual server code
	// This is where you'd call the code from cmd/cli/serve/main.go
	fmt.Println("\nTo implement: Move server logic from cmd/cli/serve/main.go here")
}
