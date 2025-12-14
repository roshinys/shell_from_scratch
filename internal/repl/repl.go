package repl

import (
	"strings"

	"github.com/roshinys/shell-from-scratch/internal/completion"
	executor "github.com/roshinys/shell-from-scratch/internal/executors"
	"github.com/roshinys/shell-from-scratch/internal/parser"
	"github.com/roshinys/shell-from-scratch/internal/redirection"
	"github.com/roshinys/shell-from-scratch/internal/terminal"
)

func Start() {
	if err := executor.InitializeExternalCommands(); err != nil {
		terminal.PrintError("initialization error: %v\n", err)
		return
	}

	for {
		DisplayPrompt()

		reader, err := completion.NewInputReader(DisplayPrompt)
		if err != nil {
			continue
		}

		fullCmd, err := reader.ReadLine()
		reader.Close()

		if err != nil {
			continue
		}

		fullCmd = strings.TrimSpace(fullCmd)
		if fullCmd == "" {
			continue
		}

		command := parser.Parse(fullCmd)

		cleanup, err := redirection.Setup(command)
		if err != nil {
			terminal.PrintError("redirection error: %v\n", err)
			continue
		}

		if err := executor.Execute(command); err != nil {
			terminal.PrintError("%s\n", err)
		}

		if cleanup != nil {
			cleanup()
		}
	}
}