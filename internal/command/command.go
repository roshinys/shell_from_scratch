package command

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/roshinys/shell-from-scratch/internal/builtin"
	"github.com/roshinys/shell-from-scratch/internal/terminal"
)


type Command struct {
	Cmd          string
	Args         []string
	Stdin      	 string // Standard os.Stdin
	Stdout       string // FileNames
	Stderr       string // FileNames
	StdinAppend  bool
	StdoutAppend bool
	StderrAppend bool
}


func (command *Command) ExecuteCommand(input io.Reader,output io.Writer) {
	if builtin.IsBuiltin(command.Cmd) {
		command.ExecuteBuiltin(input,output)
	} else {
		command.ExecuteExternalCommand(input,output)
	}
}


func(command * Command) ExecuteBuiltin(input io.Reader,output io.Writer) {
	switch command.Cmd {
	case "exit":
		command.Exit()
	case "echo":
		command.Echo(output)
	case "type":		
		command.Type(output)
	case "pwd":
		command.Pwd(output)
	case "cd":
		command.Cd()
	}
}

func (command *Command) ExecuteExternalCommand(input io.Reader,output io.Writer) {
	path, err := exec.LookPath(command.Cmd)
	if err != nil {
		terminal.PrintError("%s: command not found\n", command.Cmd)
		return
	}

	c := exec.Command(path, command.Args...)
	c.Stdin = input
	c.Stdout = output
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

func (command *Command) Echo(output io.Writer)(string, error) {
	arg := strings.Join(command.Args, " ")
    fmt.Fprintln(output, arg)
	return arg,nil
}

func (command *Command) Exit() {
	os.Exit(0)
}

func (command *Command) Type(output io.Writer) {
	if len(command.Args) == 0 {
		terminal.PrintError("type: missing argument\n")
		return
	}
	cmdToCheck := command.Args[0]
	if builtin.Builtins[cmdToCheck] {
        fmt.Fprintf(output, "%s is a shell builtin\n", cmdToCheck)
	} else if path, err := exec.LookPath(cmdToCheck); err == nil {
        fmt.Fprintf(output, "%s is %s\n", cmdToCheck, path)
	} else {
		terminal.PrintError("%s: not found\n", cmdToCheck)
	}
}

func (command *Command) Pwd(output io.Writer) {
	dir, err := os.Getwd()
	if err != nil {
		terminal.PrintError("pwd: %v\n", err)
		return 
	} 
	fmt.Fprintf(output, "%s\n", dir)
}
