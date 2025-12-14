package terminal

import "fmt"

// ANSI color codes
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorBold    = "\033[1m"
)

// ANSI control sequences
const (
	ClearLine      = "\033[2K"  // Clear entire line
	CarriageReturn = "\r"        // Move cursor to start of line
)

// Color type for reusable print functions
type Color string

const (
	Red     Color = ColorRed
	Green   Color = ColorGreen
	Yellow  Color = ColorYellow
	Blue    Color = ColorBlue
	Magenta Color = ColorMagenta
	Cyan    Color = ColorCyan
)


// Reusable print functions
func PrintColor(color Color, format string, args ...interface{}) {
	fmt.Printf("%s"+format+"%s", append([]interface{}{string(color)}, append(args, ColorReset)...)...)
}

func PrintError(format string, args ...interface{}) {
	PrintColor(Red, format, args...)
}

func PrintSuccess(format string, args ...interface{}) {
	PrintColor(Green, format, args...)
}

func PrintInfo(format string, args ...interface{}) {
	PrintColor(Blue, format, args...)
}

func PrintWarning(format string, args ...interface{}) {
	PrintColor(Yellow, format, args...)
}

func PrintPath(format string, args ...interface{}) {
	PrintColor(Magenta, format, args...)
}

// EnsureNewLine ensures we're on a fresh line at column 0
func EnsureNewLine() {
	fmt.Print("\r\n")
}

// ClearCurrentLine clears the current line and moves cursor to start
func ClearCurrentLine() {
	fmt.Print(CarriageReturn + ClearLine)
}