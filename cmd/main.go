package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/roshinys/shell-from-scratch/terminal"
)

type Command struct {
	cmd  string
	args []string
}

// Builtin commands lookup
var builtins = map[string]bool{
	"exit": true,
	"echo": true,
	"type": true,
	"pwd":  true,
	"cd":   true,
}

func main() {
	repl()
}

func repl() {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		dir,_ := os.Getwd()
		fmt.Printf("%s%s%s%s%s $ %s", terminal.ColorBold, terminal.ColorRed, getShortPath(dir), terminal.ColorReset, terminal.ColorGreen, terminal.ColorReset)
		fullCmd, err := reader.ReadString('\n')
		if err != nil {
			terminal.PrintError("Error reading input: %v\n", err)
			continue
		}

		fullCmd = strings.TrimSpace(fullCmd)
		if fullCmd == "" {
			continue
		}

		command := parseCommand(fullCmd)

		switch command.cmd {
		case "exit":
			os.Exit(0)

		case "echo":
			terminal.PrintSuccess(strings.Join(command.args, " "))

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
				terminal.PrintPath("%s\n", path)
			} else {
				terminal.PrintError("%s: not found\n", cmdToCheck)
			}

		case "pwd":
			dir, err := os.Getwd()
			if err != nil {
				terminal.PrintError("pwd: %v\n", err)
			} else {
				terminal.PrintInfo(dir)
			}

		case "cd":
			if err := changeDir(command); err != nil {
				terminal.PrintError("cd: %v\n", err)
			}

		default:
			executeCommand(command.cmd, command.args...)
		}
	}
}

func changeDir(command Command) (error){
	if len(command.args) < 1{
		return fmt.Errorf("requires atleast one dir")
	}
	path := strings.ReplaceAll(command.args[0], "~", os.Getenv("HOME"))
	err := os.Chdir(path)
	if err != nil{
		return fmt.Errorf("failed to change dir")
	}	
	return nil
}

func parseCommand(cmds string) Command {
	parts := strings.Fields(cmds)
	if len(parts) == 0 {
		return Command{}
	}

	return Command{
		cmd:  parts[0],
		args: parts[1:],
	}
}

func executeCommand(cmd string, args ...string) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Printf("%s: command not found\n", cmd)
		return
	}
	
	c := exec.Command(path, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing %s: %v\n", cmd, err)
	}
}

func getShortPath(path string) string{
	home, err := os.UserHomeDir()
	if err == nil && strings.HasPrefix(path, home) {
		return "~" + strings.TrimPrefix(path, home)
	}
	return path
}
