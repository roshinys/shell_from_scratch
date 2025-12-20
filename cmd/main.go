package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/roshinys/shell-from-scratch/internal/builtin"
	"github.com/roshinys/shell-from-scratch/internal/command"
	"github.com/roshinys/shell-from-scratch/internal/input"
)

func main() {
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
				if _, ok := builtin.Builtins[name]; !ok {
					builtin.ExternalBuiltins[name] = true
				}
			}

		}
	}
	
	return nil
}

func repl() {
	for {
		input.PrintPrompt()

		fullCmd, err := input.ReadLineWithTabCompletion()
		if err != nil {
			continue
		}
		fullCmd = strings.TrimSpace(fullCmd)
		if fullCmd == "" {
			continue
		}
		pipeline := command.ParsePipeline(fullCmd)
		pipeline.ExecutePipeline()
	}
}




