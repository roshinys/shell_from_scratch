package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/roshinys/shell-from-scratch/internal/builtin"
	"github.com/roshinys/shell-from-scratch/internal/terminal"
)


type Command struct {
	Cmd          string
	Args         []string
	Stdout       string
	Stderr       string
	StdoutAppend bool
	StderrAppend bool
}

type Pipeline struct{
	Commands []Command
}


func (command *Command) ExecuteCommand() {
	if builtin.IsBuiltin(command.Cmd) {
		command.ExecuteBuiltin()
	} else {
		command.ExecuteExternalCommand()
	}
}


func(command * Command) ExecuteBuiltin() {
	switch command.Cmd {
	case "exit":
		command.Exit()
	case "echo":
		command.Echo()
	case "type":		
		command.Type()
	case "pwd":
		command.Pwd()
	case "cd":
		command.Cd()
	}
}

func (command *Command) ExecuteExternalCommand() {
	path, err := exec.LookPath(command.Cmd)
	if err != nil {
		terminal.PrintError("%s: command not found\n", command.Cmd)
		return
	}

	c := exec.Command(path, command.Args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		terminal.PrintError("Error executing %s: %v\n", command.Cmd, err)
	}
}


func (command * Command) Cd() {
	if err := command.ChangeDir(); err != nil {
		terminal.PrintError("cd: %v\n", err)
	}
}

func (command *Command) ChangeDir() error {
	if len(command.Args) < 1 {
		return fmt.Errorf("requires atleast one dir")
	}
	path := strings.ReplaceAll(command.Args[0], "~", os.Getenv("HOME"))
	err := os.Chdir(path)
	if err != nil {
		return fmt.Errorf("failed to change dir")
	}
	return nil
}

func (command *Command) Echo()(string, error) {
	return strings.Join(command.Args, " "), nil
}

func (command *Command) Exit() {
	os.Exit(0)
}

func (command *Command) Type() {
	if len(command.Args) == 0 {
		terminal.PrintError("type: missing argument\n")
		return
	}
	cmdToCheck := command.Args[0]
	if builtin.Builtins[cmdToCheck] {
		terminal.PrintSuccess("%s", cmdToCheck)
		terminal.PrintInfo(" is a shell builtin\n")
	} else if path, err := exec.LookPath(cmdToCheck); err == nil {
		terminal.PrintSuccess("%s", cmdToCheck)
		terminal.PrintPath(" is %s\n", path)
	} else {
		terminal.PrintError("%s: not found\n", cmdToCheck)
	}
}

func (command *Command) Pwd() {
	dir, err := os.Getwd()
	if err != nil {
		terminal.PrintError("pwd: %v\n", err)
	} else {
		terminal.PrintInfo("%s\n", dir)
	}
}
