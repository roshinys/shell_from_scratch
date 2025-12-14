package executor

import (
	"os"

	"github.com/roshinys/shell-from-scratch/internal/parser"
)

// Execute dispatches commands to appropriate handlers
func Execute(command parser.Command) error {
	switch command.Cmd {
	case "exit":
		os.Exit(0)
		return nil

	case "echo":
		return ExecuteEcho(command)

	case "type":
		return ExecuteType(command, ExternalBuiltins)

	case "pwd":
		return ExecutePwd()

	case "cd":
		return ExecuteCd(command)

	default:
		return ExecuteExternal(command.Cmd, command.Args...)
	}
}

// ExecuteCommand is the old function name for backward compatibility
// It's just an alias to Execute
func ExecuteCommand(command parser.Command) error {
	return Execute(command)
}