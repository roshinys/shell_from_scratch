package parser

import "strings"

func Parse(input string) Command {
	isSingleQuote := false
	isDoubleQuote := false
	n := len(input)
	currToken := strings.Builder{}
	tokens := []string{}

	for i := 0; i < n; i++ {
		char := input[i]

		if char == '\\' && !isSingleQuote && i+1 < n {
			nextChar := input[i+1]
			if isDoubleQuote && (nextChar == '"' || nextChar == '\\' || nextChar == 'n') {
				if nextChar == 'n' {
					currToken.WriteByte('\\')
					currToken.WriteByte('n')
				} else {
					currToken.WriteByte(nextChar)
				}
				i++
				continue
			}
			if !isDoubleQuote && nextChar == ' ' {
				currToken.WriteByte(' ')
				i++
				continue
			}
		}

		if char == ' ' && !isSingleQuote && !isDoubleQuote {
			if currToken.Len() > 0 {
				tokens = append(tokens, currToken.String())
				currToken.Reset()
			}
			continue
		}

		if char == '\'' && !isDoubleQuote {
			isSingleQuote = !isSingleQuote
			continue
		}

		if char == '"' && !isSingleQuote {
			isDoubleQuote = !isDoubleQuote
			continue
		}

		currToken.WriteByte(char)
	}

	if currToken.Len() > 0 {
		tokens = append(tokens, currToken.String())
	}

	if len(tokens) == 0 {
		return Command{}
	}

	cmd := Command{Cmd: tokens[0]}
	args := []string{}
	tokenLen := len(tokens)

	for i := 1; i < tokenLen; i++ {
		token := tokens[i]
		switch token {
		case ">", "1>":
			cmd.Stdout = tokens[i+1]
			cmd.StdoutAppend = false
			i++
		case ">>", "1>>":
			cmd.Stdout = tokens[i+1]
			cmd.StdoutAppend = true
			i++
		case "2>":
			cmd.Stderr = tokens[i+1]
			cmd.StderrAppend = false
			i++
		case "2>>":
			cmd.Stderr = tokens[i+1]
			cmd.StderrAppend = true
			i++
		default:
			args = append(args, token)
		}
	}
	cmd.Args = args
	return cmd
}