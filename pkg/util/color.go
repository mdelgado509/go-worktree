// Package util provides utility functions
package util

// ANSI color codes for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// Colorize returns a string with color codes
func Colorize(text, color string) string {
	return color + text + ColorReset
}

// Bold returns a string in bold
func Bold(text string) string {
	return "\033[1m" + text + "\033[0m"
}
