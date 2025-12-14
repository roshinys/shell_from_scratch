package repl

import (
	"fmt"
	"os"
	"strings"

	"github.com/roshinys/shell-from-scratch/internal/terminal"
)

func DisplayPrompt() {
	dir, _ := os.Getwd()
	fmt.Printf("%s%s%s%s%s$ %s",
		terminal.ColorBold, terminal.ColorCyan,
		getShortPath(dir), terminal.ColorReset,
		terminal.ColorGreen, terminal.ColorReset)
}

func getShortPath(path string) string {
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME")
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	home, err := os.UserHomeDir()
	shortPath := path
	if err == nil && strings.HasPrefix(path, home) {
		shortPath = "~" + strings.TrimPrefix(path, home)
	}

	return username + "@" + hostname + ":" + shortPath
}