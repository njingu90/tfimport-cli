package cmd

import (
	"fmt"
	"os"
	"strings"
)

// Color codes for console output
const (
	ColorReset   = "\033[0m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorRed     = "\033[31m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
)

// PrintSuccess prints a success message
func PrintSuccess(msg string) {
	if supportsColor() {
		fmt.Printf("%s✓%s %s\n", ColorGreen, ColorReset, msg)
	} else {
		fmt.Printf("[OK] %s\n", msg)
	}
}

// PrintError prints an error message
func PrintError(msg string) {
	if supportsColor() {
		fmt.Fprintf(os.Stderr, "%s✗%s %s\n", ColorRed, ColorReset, msg)
	} else {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", msg)
	}
}

// PrintWarning prints a warning message
func PrintWarning(msg string) {
	if supportsColor() {
		fmt.Fprintf(os.Stderr, "%s⚠%s %s\n", ColorYellow, ColorReset, msg)
	} else {
		fmt.Fprintf(os.Stderr, "[WARN] %s\n", msg)
	}
}

// PrintInfo prints an info message
func PrintInfo(msg string) {
	if supportsColor() {
		fmt.Printf("%sℹ%s %s\n", ColorBlue, ColorReset, msg)
	} else {
		fmt.Printf("[INFO] %s\n", msg)
	}
}

// PrintSection prints a section header
func PrintSection(title string) {
	if supportsColor() {
		fmt.Printf("\n%s=== %s ===%s\n", ColorMagenta, title, ColorReset)
	} else {
		fmt.Printf("\n=== %s ===\n", title)
	}
}

// PrintKeyValue prints a key-value pair
func PrintKeyValue(key string, value interface{}) {
	fmt.Printf("  %-20s %v\n", key+":", value)
}

// PrintList prints a list of items
func PrintList(items []string) {
	for _, item := range items {
		fmt.Printf("  - %s\n", item)
	}
}

// PrintTable prints a simple table
func PrintTable(headers []string, rows [][]string) {
	// Calculate column widths
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print header
	for i, header := range headers {
		fmt.Printf("%-*s  ", widths[i], header)
	}
	fmt.Println()

	// Print separator
	for i, w := range widths {
		fmt.Print(strings.Repeat("-", w))
		if i < len(widths)-1 {
			fmt.Print("  ")
		}
	}
	fmt.Println()

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				fmt.Printf("%-*s  ", widths[i], cell)
			}
		}
		fmt.Println()
	}
}

// supportsColor checks if the terminal supports color output
func supportsColor() bool {
	// Check for NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check for TERM environment variable
	term := os.Getenv("TERM")
	if term == "dumb" {
		return false
	}

	return true
}

// ProgressBar represents a simple progress indicator
type ProgressBar struct {
	current int
	total   int
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{total: total}
}

// Increment increases the progress
func (pb *ProgressBar) Increment() {
	pb.current++
	if pb.current > pb.total {
		pb.current = pb.total
	}
}

// Print prints the progress bar
func (pb *ProgressBar) Print() {
	if pb.total == 0 {
		return
	}

	percent := (pb.current * 100) / pb.total
	filled := (pb.current * 20) / pb.total

	fmt.Printf("\r[%-20s] %d%% (%d/%d)", strings.Repeat("=", filled)+strings.Repeat(" ", 20-filled), percent, pb.current, pb.total)
}

// Done marks the progress bar as complete
func (pb *ProgressBar) Done() {
	fmt.Println()
}
