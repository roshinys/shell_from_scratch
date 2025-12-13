package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)


type Terminal struct{
	id int 
}


type TerminalManager struct{
	terminals map[int]*Terminal
	nextId int
	mu sync.RWMutex
}


type Command struct{
	cmd string 
	args []string
}



func NewTerminalManager() *TerminalManager{
	return &TerminalManager{
		terminals: make(map[int]*Terminal),
		nextId: 1,
	}
}

func (tm *TerminalManager) CreateTerminal() *Terminal{
	tm.mu.Lock()
	defer tm.mu.Unlock()

	t := &Terminal{
		id: tm.nextId,
	}
	tm.terminals[tm.nextId] = t;
	tm.nextId++
	return t
}


func (tm *TerminalManager) GetTerminal(id int) (*Terminal, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	t, exists := tm.terminals[id]
	return t, exists
}

func (tm *TerminalManager) CloseTerminal(id int) bool {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	if _, exists := tm.terminals[id]; exists {
		delete(tm.terminals, id)
		return true
	}
	return false
}

func (tm *TerminalManager) ListTerminals() []int {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	ids := make([]int, 0, len(tm.terminals))
	for id := range tm.terminals {
		ids = append(ids, id)
	}
	return ids
}


func (tm *TerminalManager) RunTerminal(t *Terminal) {
	for {
		fmt.Print("$ ")
		fullCmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}

		fullCmd = strings.TrimSpace(fullCmd)
		if fullCmd == "" {
			continue
		}

		command := parseCommand(fullCmd)

		switch command.cmd {
		case "exit":
			os.Exit(1)

		case "echo":
			fmt.Println(strings.Join(command.args, " "))

		default:
			fmt.Println(command.cmd + ": command not found")
		}
	}
}


func parseCommand(cmds string)Command{
	cmds = strings.TrimSpace(cmds)
	parts := strings.Fields(cmds)
	if len(parts) == 0 {
		return Command{}
	}

	cmd := parts[0]
	args := parts[1:]
	return Command{
		cmd: cmd,
		args: args,
	}
}
