package terminal

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
