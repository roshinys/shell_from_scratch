package completion

import (
	"fmt"
	"os"

	executor "github.com/roshinys/shell-from-scratch/internal/executors"
	"golang.org/x/term"
)

type InputReader struct {
	fd          int
	oldState    *term.State
	displayFunc func()
}

func NewInputReader(displayFunc func()) (*InputReader, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}

	return &InputReader{
		fd:          fd,
		oldState:    oldState,
		displayFunc: displayFunc,
	}, nil
}

func (r *InputReader) Close() {
	term.Restore(r.fd, r.oldState)
}

func (r *InputReader) ReadLine() (string, error) {
	buf := make([]byte, 1)
	currCmd := ""
	lastWasTab := false

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			return "", err
		}

		char := buf[0]

		// ENTER
		if char == '\r' || char == '\n' {
			fmt.Print("\n\r")
			return currCmd, nil
		}

		// CTRL+C
		if char == 3 {
			fmt.Print("\n\r")
			return "", fmt.Errorf("interrupted")
		}

		// BACKSPACE
		if char == 127 {
			if len(currCmd) > 0 {
				currCmd = currCmd[:len(currCmd)-1]
				fmt.Print("\b \b")
			}
			lastWasTab = false
			continue
		}

		// TAB
		if char == '\t' {
			completions := GetCompletions(currCmd, executor.Builtins, executor.ExternalBuiltins)

			if len(completions) == 0 {
				lastWasTab = false
				continue
			}

			// Single match
			if len(completions) == 1 {
				match := completions[0]
				remaining := match[len(currCmd):]
				fmt.Print(remaining + " ")
				currCmd = match + " "
				lastWasTab = false
				continue
			}

			// Multiple matches
			lcp := LongestCommonPrefix(completions)

			if len(lcp) > len(currCmd) {
				remaining := lcp[len(currCmd):]
				fmt.Print(remaining)
				currCmd = lcp
				lastWasTab = false
			} else if lastWasTab {
				// Second TAB
				fmt.Print("\n\r")
				for _, m := range completions {
					fmt.Print(m + "  ")
				}
				fmt.Print("\n\r")
				r.displayFunc()
				fmt.Print(currCmd)
				lastWasTab = false
			} else {
				lastWasTab = true
			}

			continue
		}

		// NORMAL CHAR
		if char >= 32 && char < 127 {
			currCmd += string(char)
			fmt.Print(string(char))
			lastWasTab = false
		}
	}
}