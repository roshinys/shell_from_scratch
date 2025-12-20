// internal/terminal/terminal.go
package terminal

import "fmt"

// Color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

// PrintError prints error messages in red
func PrintError(format string, args ...interface{}) {
	fmt.Printf(ColorRed+format+ColorReset, args...)
}

// PrintSuccess prints success messages in green
func PrintSuccess(format string, args ...interface{}) {
	fmt.Printf(ColorGreen+format+ColorReset, args...)
}

// PrintInfo prints info messages in cyan
func PrintInfo(format string, args ...interface{}) {
	fmt.Printf(ColorCyan+format+ColorReset, args...)
}

// PrintPath prints paths in blue
func PrintPath(format string, args ...interface{}) {
	fmt.Printf(ColorBlue+format+ColorReset, args...)
}

// PrintWarning prints warning messages in yellow
func PrintWarning(format string, args ...interface{}) {
	fmt.Printf(ColorYellow+format+ColorReset, args...)
}