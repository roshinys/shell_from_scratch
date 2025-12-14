package executor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var ExternalBuiltins = make(map[string]bool)

func InitializeExternalCommands() error {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return fmt.Errorf("reading path failed")
	}

	pathDirs := filepath.SplitList(pathEnv)
	for _, dir := range pathDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			name := entry.Name()
			if info.Mode().IsRegular() && info.Mode()&0111 != 0 {
				if _, ok := Builtins[name]; !ok {
					ExternalBuiltins[name] = true
				}
			}
		}
	}
	return nil
}

func ExecuteExternal(cmd string, args ...string) error {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return fmt.Errorf("%s: command not found", cmd)
	}

	c := exec.Command(path, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return fmt.Errorf("error executing %s: %v", cmd, err)
	}
	return nil
}