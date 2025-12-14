package redirection

import (
	"os"

	"github.com/roshinys/shell-from-scratch/internal/parser"
)

func Setup(command parser.Command) (func(), error) {
	var oldStdout *os.File
	var oldStderr *os.File

	if command.Stdout != "" {
		oldStdout = os.Stdout
		var f *os.File
		var err error
		if command.StdoutAppend {
			f, err = os.OpenFile(command.Stdout, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		} else {
			f, err = os.Create(command.Stdout)
		}
		if err != nil {
			return nil, err
		}
		os.Stdout = f
	}

	if command.Stderr != "" {
		oldStderr = os.Stderr
		var f *os.File
		var err error
		if command.StderrAppend {
			f, err = os.OpenFile(command.Stderr, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		} else {
			f, err = os.Create(command.Stderr)
		}
		if err != nil {
			return nil, err
		}
		os.Stderr = f
	}

	cleanup := func() {
		if oldStdout != nil {
			os.Stdout.Close()
			os.Stdout = oldStdout
		}
		if oldStderr != nil {
			os.Stderr.Close()
			os.Stderr = oldStderr
		}
	}

	return cleanup, nil
}