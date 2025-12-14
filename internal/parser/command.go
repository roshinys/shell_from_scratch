package parser

type Command struct {
	Cmd          string
	Args         []string
	Stdout       string
	Stderr       string
	StdoutAppend bool
	StderrAppend bool
}