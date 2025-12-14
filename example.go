package shellfromscratch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/roshinys/shell-from-scratch/internal/terminal"
	"golang.org/x/term"
)

type Command struct {
	cmd          string
	args         []string
	stdout       string
	stderr       string
	stdoutAppend bool
	stderrAppend bool
}

// Builtin commands lookup
var builtins = map[string]bool{
	"exit": true,
	"echo": true,
	"type": true,
	"pwd":  true,
	"cd":   true,
}

var externalBuiltins = make(map[string]bool)

func exampleMain() {
	initRepl()
	repl()
}

func initRepl() (error){
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return fmt.Errorf("reading path failed")
	}
	pathDirs := filepath.SplitList(pathEnv)
	for _,dir := range pathDirs{
		entries,err := os.ReadDir(dir)
		if err != nil{
			continue
		}
		for _,entry := range entries{
			info,err := entry.Info()
			if err != nil {
    			continue
			}
			// check for executable and has permission to execute the function
			name := entry.Name()
			if info.Mode().IsRegular() && info.Mode()&0111 != 0 {
				if _, ok := builtins[name]; !ok {
					externalBuiltins[name] = true
				}
			}

		}
	}
	
	return nil
}

func repl() {
	for {
		printPrompt()

		fullCmd, err := readLineWithTabCompletion()
		if err != nil {
			continue
		}

		fullCmd = strings.TrimSpace(fullCmd)
		if fullCmd == "" {
			continue
		}

		command := parseCommand(fullCmd)
		cleanup, err := setupRedirection(command)
		if err != nil {
			terminal.PrintError("redirection error: %v\n", err)
			continue
		}

		switch command.cmd {
		case "exit":
			os.Exit(0)

		case "echo":
			arg, err := echoPrint(command)
			if err != nil {
				terminal.PrintError("%s\n", err)
				continue
			}
			fmt.Print(arg)
			fmt.Print("\n")

		case "type":
			if len(command.args) == 0 {
				terminal.PrintError("type: missing argument\n")
				continue
			}
			cmdToCheck := command.args[0]
			if builtins[cmdToCheck] {
				terminal.PrintSuccess("%s", cmdToCheck)
				terminal.PrintInfo(" is a shell builtin\n")
			} else if path, err := exec.LookPath(cmdToCheck); err == nil {
				terminal.PrintSuccess("%s", cmdToCheck)
				terminal.PrintPath(" is %s\n", path)
			} else {
				terminal.PrintError("%s: not found\n", cmdToCheck)
			}

		case "pwd":
			dir, err := os.Getwd()
			if err != nil {
				terminal.PrintError("pwd: %v\n", err)
			} else {
				terminal.PrintInfo("%s\n", dir)
			}

		case "cd":
			if err := changeDir(command); err != nil {
				terminal.PrintError("cd: %v\n", err)
			}

		default:
			executeCommand(command.cmd, command.args...)
		}
		if cleanup != nil {
			cleanup()
		}
	}
}

func printPrompt() {
	dir, _ := os.Getwd()
	fmt.Printf("%s%s%s%s%s$ %s",
		terminal.ColorBold, terminal.ColorCyan,
		getShortPath(dir), terminal.ColorReset,
		terminal.ColorGreen, terminal.ColorReset)
}

func readLineWithTabCompletion() (string, error) {
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
			completions := getCompletions(currCmd)

			if len(completions) == 0 {
				lastWasTab = false
				continue
			}

			// Single match → complete fully
			if len(completions) == 1 {
				match := completions[0]
				remaining := match[len(currCmd):]
				fmt.Print(remaining + " ")
				currCmd = match + " "
				lastWasTab = false
				continue
			}

			// Multiple matches
			lcp := longestCommonPrefix(completions)

			if len(lcp) > len(currCmd) {
				// Extend to common prefix
				remaining := lcp[len(currCmd):]
				fmt.Print(remaining)
				currCmd = lcp
				lastWasTab = false
			} else if lastWasTab {
				// Second TAB → show all options
				fmt.Print("\n\r")
				for _, m := range completions {
					fmt.Print(m + "  ")
				}
				fmt.Print("\n\r")
				printPrompt()
				fmt.Print(currCmd)
				lastWasTab = false
			} else {
				// First TAB, no progress
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


func getCompletions(prefix string) []string {
	matches := []string{}

	// Only autocomplete "echo" and "exit" as per requirements
	builtinsInternal := mapKeys(builtins)
	builtinsExternal := mapKeys(externalBuiltins)

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

func echoPrint(command Command) (string, error) {
	return strings.Join(command.args, " "), nil
}

func changeDir(command Command) error {
	if len(command.args) < 1 {
		return fmt.Errorf("requires atleast one dir")
	}
	path := strings.ReplaceAll(command.args[0], "~", os.Getenv("HOME"))
	err := os.Chdir(path)
	if err != nil {
		return fmt.Errorf("failed to change dir")
	}
	return nil
}

func parseCommand(input string) Command {
	isSingleQuote := false
	isDoubleQuote := false
	n := len(input)
	currToken := strings.Builder{}
	tokens := []string{}

	for i := 0; i < n; i++ {
		char := input[i]

		// Handle escape sequences in double quotes
		if char == '\\' && !isSingleQuote && i+1 < n {
			nextChar := input[i+1]
			// In double quotes, only certain chars are escaped
			if isDoubleQuote && (nextChar == '"' || nextChar == '\\' || nextChar == 'n') {
				if nextChar == 'n' {
					currToken.WriteByte('\\')
					currToken.WriteByte('n')
				} else {
					currToken.WriteByte(nextChar)
				}
				i++ // skip next char
				continue
			}
			// Outside quotes, backslash escapes spaces
			if !isDoubleQuote && nextChar == ' ' {
				currToken.WriteByte(' ')
				i++
				continue
			}
		}

		if char == ' ' && !isSingleQuote && !isDoubleQuote {
			if currToken.Len() > 0 {
				tokens = append(tokens, currToken.String())
				currToken.Reset()
			}
			continue
		}

		if char == '\'' && !isDoubleQuote {
			isSingleQuote = !isSingleQuote
			continue
		}

		if char == '"' && !isSingleQuote {
			isDoubleQuote = !isDoubleQuote
			continue
		}

		currToken.WriteByte(char)
	}

	if currToken.Len() > 0 {
		tokens = append(tokens, currToken.String())
	}

	if len(tokens) == 0 {
		return Command{}
	}
	cmd := Command{
		cmd: tokens[0],
	}
	args := []string{}
	tokenLen := len(tokens)
	for i := 1; i < tokenLen; i++ {  // Start from 1, not 0!
		token := tokens[i]
		switch token {
		case ">", "1>":
			cmd.stdout = tokens[i+1]
			cmd.stdoutAppend = false
			i++
		case ">>", "1>>":
			cmd.stdout = tokens[i+1]
			cmd.stdoutAppend = true
			i++

		case "2>":
			cmd.stderr = tokens[i+1]
			cmd.stderrAppend = false
			i++

		case "2>>":
			cmd.stderr = tokens[i+1]
			cmd.stderrAppend = true
			i++

		default:
			args = append(args, token)
		}
	}
	cmd.args = args
	return cmd
}

func setupRedirection(command Command) (func(), error) {
	var oldStdout *os.File
	var oldStderr *os.File

	if command.stdout != "" {
		oldStdout = os.Stdout
		var f *os.File
		var err error
		if command.stdoutAppend {
			f, err = os.OpenFile(command.stdout, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		} else {
			f, err = os.Create(command.stdout)
		}
		if err != nil {
			return nil, err
		}
		os.Stdout = f
	}

	if command.stderr != "" {
		oldStderr = os.Stderr
		var f *os.File
		var err error
		if command.stderrAppend {
			f, err = os.OpenFile(command.stderr, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		} else {
			f, err = os.Create(command.stderr)
		}
		if err != nil {
			return nil, err
		}
		os.Stderr = f
	}

	return func() {
		if oldStdout != nil {
			os.Stdout.Close()
			os.Stdout = oldStdout
		}
		if oldStderr != nil {
			os.Stderr.Close()
			os.Stderr = oldStderr
		}
	}, nil
}

func executeCommand(cmd string, args ...string) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		terminal.PrintError("%s: command not found\n", cmd)
		return
	}

	c := exec.Command(path, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		terminal.PrintError("Error executing %s: %v\n", cmd, err)
	}
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