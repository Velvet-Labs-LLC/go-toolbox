// Package cli provides utilities for building command-line interfaces.
package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// OutputFormat represents different output formats.
type OutputFormat string

const (
	OutputTable OutputFormat = "table"
	OutputJSON  OutputFormat = "json"
	OutputYAML  OutputFormat = "yaml"
)

// Colors for different message types.
var (
	InfoColor    = color.New(color.FgCyan)
	SuccessColor = color.New(color.FgGreen)
	WarnColor    = color.New(color.FgYellow)
	ErrorColor   = color.New(color.FgRed)
	HeaderColor  = color.New(color.FgBlue, color.Bold)
)

// BaseCommand provides common functionality for CLI commands.
type BaseCommand struct {
	*cobra.Command
	Verbose bool
	Output  OutputFormat
}

// NewBaseCommand creates a new base command with common flags.
func NewBaseCommand(use, short string) *BaseCommand {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
	}

	baseCmd := &BaseCommand{
		Command: cmd,
	}

	// Add common flags
	cmd.PersistentFlags().BoolVarP(&baseCmd.Verbose, "verbose", "v", false, "Enable verbose output")
	cmd.PersistentFlags().StringVar((*string)(&baseCmd.Output), "output", "table", "Output format (table, json, yaml)")

	return baseCmd
}

// PrintInfo prints an info message.
func (c *BaseCommand) PrintInfo(format string, args ...interface{}) {
	if c.Output == OutputTable {
		InfoColor.Printf(format+"\n", args...)
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

// PrintSuccess prints a success message.
func (c *BaseCommand) PrintSuccess(format string, args ...interface{}) {
	if c.Output == OutputTable {
		SuccessColor.Printf(format+"\n", args...)
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

// PrintWarn prints a warning message.
func (c *BaseCommand) PrintWarn(format string, args ...interface{}) {
	if c.Output == OutputTable {
		WarnColor.Printf(format+"\n", args...)
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

// PrintError prints an error message.
func (c *BaseCommand) PrintError(format string, args ...interface{}) {
	if c.Output == OutputTable {
		ErrorColor.Printf(format+"\n", args...)
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

// PrintHeader prints a header message.
func (c *BaseCommand) PrintHeader(format string, args ...interface{}) {
	if c.Output == OutputTable {
		HeaderColor.Printf(format+"\n", args...)
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

// PrintVerbose prints a message only if verbose mode is enabled.
func (c *BaseCommand) PrintVerbose(format string, args ...interface{}) {
	if c.Verbose {
		c.PrintInfo(format, args...)
	}
}

// Table provides utilities for creating tables.
type Table struct {
	writer  *tablewriter.Table
	data    [][]string
	headers []string
}

// NewTable creates a new table.
func NewTable(headers []string) *Table {
	table := tablewriter.NewWriter(os.Stdout)

	return &Table{
		writer:  table,
		data:    make([][]string, 0),
		headers: headers,
	}
}

// AddRow adds a row to the table.
func (t *Table) AddRow(row ...string) {
	t.data = append(t.data, row)
}

// Render renders the table.
func (t *Table) Render() {
	// Set headers if they exist
	if len(t.headers) > 0 {
		t.writer.Header(toInterfaceSlice(t.headers)...)
	}

	// Add all data rows
	for _, row := range t.data {
		t.writer.Append(toInterfaceSlice(row)...)
	}

	// Render the table
	t.writer.Render()
}

// toInterfaceSlice converts []string to []interface{}.
func toInterfaceSlice(strings []string) []interface{} {
	interfaces := make([]interface{}, len(strings))
	for i, s := range strings {
		interfaces[i] = s
	}
	return interfaces
}

// ProgressBar creates a progress bar.
type ProgressBar struct {
	bar *progressbar.ProgressBar
}

// NewProgressBar creates a new progress bar.
func NewProgressBar(maxValue int, description string) *ProgressBar {
	bar := progressbar.NewOptions(maxValue,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWidth(50),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)

	return &ProgressBar{bar: bar}
}

// Add increments the progress bar.
func (p *ProgressBar) Add(num int) {
	p.bar.Add(num)
}

// Finish completes the progress bar.
func (p *ProgressBar) Finish() {
	p.bar.Finish()
}

// Prompt provides utilities for user input.
type Prompt struct{}

// NewPrompt creates a new prompt.
func NewPrompt() *Prompt {
	return &Prompt{}
}

// String prompts for a string input.
func (p *Prompt) String(label string, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

// Password prompts for a password input.
func (p *Prompt) Password(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

// Confirm prompts for a yes/no confirmation.
func (p *Prompt) Confirm(label string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label + " (y/N)",
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}

	return strings.ToLower(result) == "y" || strings.ToLower(result) == "yes", nil
}

// Select prompts for selection from a list.
func (p *Prompt) Select(label string, items []string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	index, result, err := prompt.Run()
	if err != nil {
		return -1, "", err
	}

	return index, result, nil
}

// Spinner provides a simple spinner for long-running operations.
type Spinner struct {
	chars   []string
	current int
	stop    chan bool
}

// NewSpinner creates a new spinner.
func NewSpinner() *Spinner {
	return &Spinner{
		chars:   []string{"|", "/", "-", "\\"},
		current: 0,
		stop:    make(chan bool),
	}
}

// Start starts the spinner.
func (s *Spinner) Start(message string) {
	go func() {
		for {
			select {
			case <-s.stop:
				return
			default:
				fmt.Printf("\r%s %s", s.chars[s.current], message)
				s.current = (s.current + 1) % len(s.chars)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner.
func (s *Spinner) Stop() {
	s.stop <- true
	fmt.Print("\r")
}

// ParseSize parses a size string (e.g., "1KB", "2MB") into bytes.
func ParseSize(sizeStr string) (int64, error) {
	sizeStr = strings.ToUpper(strings.TrimSpace(sizeStr))

	multipliers := map[string]int64{
		"B":  1,
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
	}

	for suffix, multiplier := range multipliers {
		if strings.HasSuffix(sizeStr, suffix) {
			numStr := strings.TrimSuffix(sizeStr, suffix)
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid size format: %s", sizeStr)
			}
			return int64(num * float64(multiplier)), nil
		}
	}

	// If no suffix, assume bytes
	num, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size format: %s", sizeStr)
	}

	return num, nil
}

// FormatSize formats bytes into a human-readable string.
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDuration formats a duration into a human-readable string.
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.0fms", d.Seconds()*1000)
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}
