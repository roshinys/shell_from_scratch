package shell

import (
	"fmt"
	"io"
)

type Shell struct{
	History []string
}

func NewShell() *Shell{
	return &Shell{
		History: make([]string, 0),
	}
}

func (shell *Shell) AddHistory(fullCmd string) *Shell{
	shell.History = append(shell.History, fullCmd)
	return shell
}

func (shell *Shell) GetHistory(w io.Writer,k int){
	history := shell.History
	n := len(shell.History) // 2 
	if k < 0 || k > n {
		k = n
	} 
	m := n - k
	if m < 0 {
		m = 0
	}
	// get last n elements 
	for i:=n-1;i>=m;i--{
		fmt.Fprintf(w, "%d. %s\n",i+1,history[i])
	}
}
