package command

import (
	"io"
	"os"
)

type Pipeline struct{
	Commands []Command
	isPipeline bool
}

func (pipeline *Pipeline) ExecutePipeline(){
	// Assume only two commands exists (max)
	n := len(pipeline.Commands)
	if n == 1{
		pipeline.Commands[0].ExecuteCommand(os.Stdin,os.Stdout)
		return
	}
	var input io.Reader = os.Stdin
	for i:=0;i<n;i+=1{
		if i == n-1{
			pipeline.Commands[i].ExecuteCommand(input,os.Stdout)
		}else{
			in := input // Need to initialize here since it is used in go functions and input is modifying by the time go executes input would have modified
			reader, writer := io.Pipe()
			cmd := pipeline.Commands[i]
			go func() {
				defer writer.Close()  // IMPORTANT: Close after cmd finishes
				cmd.ExecuteCommand(in,writer)
			}()
			input = reader
		}
	}
}