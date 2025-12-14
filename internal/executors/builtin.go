package executor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/roshinys/shell-from-scratch/internal/parser"
	"github.com/roshinys/shell-from-scratch/internal/terminal"
)

var Builtins = map[string]bool{
	"exit": true,
	"echo": true,
	"type": true,
	"pwd":  true,
	"cd":   true,
}

func ExecuteEcho(command parser.Command) error {
	arg := strings.Join(command.Args, " ")
	fmt.Print(arg)
	fmt.Print("\n")
	return nil
}

func ExecuteType(command parser.Command, externalBuiltins map[string]bool) error {
	if len(command.Args) == 0 {
		return fmt.Errorf("type: missing argument")
	}

	cmdToCheck := command.Args[0]
	if Builtins[cmdToCheck] {
		terminal.PrintSuccess("%s", cmdToCheck)
		terminal.PrintInfo(" is a shell builtin\n")
	} else if path, err := exec.LookPath(cmdToCheck); err == nil {
		terminal.PrintSuccess("%s", cmdToCheck)
		terminal.PrintPath(" is %s\n", path)
	} else {
		return fmt.Errorf("%s: not found", cmdToCheck)
	}
	return nil
}

func ExecutePwd() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("pwd: %v", err)
	}
	terminal.PrintInfo("%s\n", dir)
	return nil
}

func ExecuteCd(command parser.Command) error {
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