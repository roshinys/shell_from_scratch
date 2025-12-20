package redirection

import (
	"os"

	"github.com/roshinys/shell-from-scratch/internal/shell"
)

func SetupRedirection(command shell.Command) (func(), error) {
	var oldStdout *os.File
	var oldStderr *os.File

	if command.Stdout != "" {
		var err error
		oldStdout, err = redirectStdout(command)
		if err != nil {
			return nil, err
		}
	}

	if command.Stderr != "" {
		var err error
		oldStderr, err = redirectStderr(command)
		if err != nil {
			return nil, err
		}
	}

	return createCleanupFunction(oldStdout, oldStderr), nil
}

func redirectStdout(command shell.Command) (*os.File, error) {
	oldStdout := os.Stdout
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
	return oldStdout, nil
}

func redirectStderr(command shell.Command) (*os.File, error) {
	oldStderr := os.Stderr
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
	return oldStderr, nil
}

func createCleanupFunction(oldStdout *os.File, oldStderr *os.File) func() {
	return func() {
		if oldStdout != nil {
			os.Stdout.Close()
			os.Stdout = oldStdout
		}
		if oldStderr != nil {
			os.Stderr.Close()
			os.Stderr = oldStderr
		}
	}
}


