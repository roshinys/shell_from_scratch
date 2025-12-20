package input

import (
	"fmt"
	"os"
	"strings"

	"github.com/roshinys/shell-from-scratch/internal/builtin"
	"github.com/roshinys/shell-from-scratch/internal/terminal"
	"golang.org/x/term"
)

func ReadLineWithTabCompletion() (string, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer term.Restore(fd, oldState)

	buf := make([]byte, 1)
	currCmd := ""
	lastWasTab := false

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			return "", err
		}

		char := buf[0]

		if isEnterKey(char) {
			fmt.Print("\n\r")
			return currCmd, nil
		}

		if isInterruptKey(char) {
			fmt.Print("\n\r")
			return "", fmt.Errorf("interrupted")
		}

		if isBackspaceKey(char) {
			currCmd = handleBackspace(currCmd)
			lastWasTab = false
			continue
		}

		if isTabKey(char) {
			currCmd, lastWasTab = handleTabCompletion(currCmd, lastWasTab)
			continue
		}

		if isPrintableChar(char) {
			currCmd += string(char)
			fmt.Print(string(char))
			lastWasTab = false
		}
	}
}


func isEnterKey(char byte) bool {
	return char == '\r' || char == '\n'
}

func isInterruptKey(char byte) bool {
	return char == 3
}

func isBackspaceKey(char byte) bool {
	return char == 127
}

func isTabKey(char byte) bool {
	return char == '\t'
}

func isPrintableChar(char byte) bool {
	return char >= 32 && char < 127
}

func handleBackspace(currCmd string) string {
	if len(currCmd) > 0 {
		currCmd = currCmd[:len(currCmd)-1]
		fmt.Print("\b \b")
	}
	return currCmd
}

func handleTabCompletion(currCmd string, lastWasTab bool) (string, bool) {
	completions := getCompletions(currCmd)

	if len(completions) == 0 {
		return currCmd, false
	}

	if len(completions) == 1 {
		return handleSingleCompletion(currCmd, completions[0])
	}

	return handleMultipleCompletions(currCmd, completions, lastWasTab)
}

func handleSingleCompletion(currCmd string, match string) (string, bool) {
	remaining := match[len(currCmd):]
	fmt.Print(remaining + " ")
	return match + " ", false
}

func handleMultipleCompletions(currCmd string, completions []string, lastWasTab bool) (string, bool) {
	lcp := longestCommonPrefix(completions)

	if len(lcp) > len(currCmd) {
		remaining := lcp[len(currCmd):]
		fmt.Print(remaining)
		return lcp, false
	}

	if lastWasTab {
		showCompletions(completions)
		PrintPrompt()
		fmt.Print(currCmd)
		return currCmd, false
	}

	return currCmd, true
}

func showCompletions(completions []string) {
	fmt.Print("\n\r")
	for _, m := range completions {
		fmt.Print(m + "  ")
	}
	fmt.Print("\n\r")
}



func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	prefix := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, prefix) {
			if len(prefix) == 0 {
				return ""
			}
			prefix = prefix[:len(prefix)-1]
		}
	}
	return prefix
}


func PrintPrompt() {
	dir, _ := os.Getwd()
	fmt.Printf("%s%s%s%s%s$ %s",
		terminal.ColorBold, terminal.ColorCyan,
		getShortPath(dir), terminal.ColorReset,
		terminal.ColorGreen, terminal.ColorReset)
}


func getShortPath(path string) string {
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME") // Windows fallback
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

func getCompletions(prefix string) []string {
	matches := []string{}

	// Only autocomplete "echo" and "exit" as per requirements
	builtinsInternal := mapKeys(builtin.Builtins)
	builtinsExternal := mapKeys(builtin.ExternalBuiltins)

	builtinsToComplete := append(builtinsInternal,builtinsExternal...)

	for _, cmd := range builtinsToComplete {
		if strings.HasPrefix(cmd, prefix) {
			matches = append(matches, cmd)
		}
	}
	return matches
}

func mapKeys(builtins map[string]bool)([]string){
	tmp := []string{}
	for k,_ := range builtins{
		tmp = append(tmp, k)
	}
	return tmp
}