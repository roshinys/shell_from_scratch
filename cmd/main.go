package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/roshinys/shell-from-scratch/internal/builtin"
	"github.com/roshinys/shell-from-scratch/internal/input"
	"github.com/roshinys/shell-from-scratch/internal/shell"
)

func main() {
	s := shell.NewShell()
	initRepl()
	repl(s)
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

func repl(s *shell.Shell) {
	for {
		input.PrintPrompt()
		fullCmd, err := input.ReadLineWithTabCompletion(s)
		if err != nil {
			continue
		}
		fullCmd = strings.TrimSpace(fullCmd)
		if fullCmd == "" {
			continue
		}
		re := regexp.MustCompile(`^history`)
		if !re.MatchString(fullCmd){
			s = s.AddHistory(fullCmd)
		}
		pipeline := shell.ParsePipeline(fullCmd)
		pipeline.ExecutePipeline(s)
	}
}




