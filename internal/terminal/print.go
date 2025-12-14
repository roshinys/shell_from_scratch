package terminal

import "fmt"

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

func EnsureNewLine() {
	fmt.Print("\r\n")
}

func ClearCurrentLine() {
	fmt.Print(CarriageReturn + ClearLine)
}