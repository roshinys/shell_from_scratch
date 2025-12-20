package command

import "strings"


func ParsePipeline(input string) Pipeline {
	inputParts := strings.Split(input,"|")
	count := len(inputParts)
	commands := []Command{}
	for i := 0; i < count; i++ {
		commands = append(commands, ParseCommand(inputParts[i]))
	}
	return Pipeline{
		Commands: commands,
	}
}

func ParseCommand(input string) Command {
	isSingleQuote := false
	isDoubleQuote := false
	n := len(input)
	currToken := strings.Builder{}
	tokens := []string{}

	for i := 0; i < n; i++ {
		char := input[i]

		if shouldHandleEscape(char, isSingleQuote, isDoubleQuote, i, n, input) {
			i = handleEscapeSequence(char, isDoubleQuote, i, n, input, &currToken)
			continue
		}

		if isUnquotedSpace(char, isSingleQuote, isDoubleQuote) {
			if currToken.Len() > 0 {
				tokens = append(tokens, currToken.String())
				currToken.Reset()
			}
			continue
		}

		if isSingleQuoteToggle(char, isDoubleQuote) {
			isSingleQuote = !isSingleQuote
			continue
		}

		if isDoubleQuoteToggle(char, isSingleQuote) {
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
	return buildCommandFromTokens(tokens)
}


func shouldHandleEscape(char byte, isSingleQuote bool, isDoubleQuote bool, i int, n int, input string) bool {
	return char == '\\' && !isSingleQuote && i+1 < n
}

func handleEscapeSequence(char byte, isDoubleQuote bool, i int, n int, input string, currToken *strings.Builder) int {
	nextChar := input[i+1]
	if isDoubleQuote && (nextChar == '"' || nextChar == '\\' || nextChar == 'n') {
		if nextChar == 'n' {
			currToken.WriteByte('\\')
			currToken.WriteByte('n')
		} else {
			currToken.WriteByte(nextChar)
		}
		return i + 1
	}
	if !isDoubleQuote && nextChar == ' ' {
		currToken.WriteByte(' ')
		return i + 1
	}
	return i
}

func isUnquotedSpace(char byte, isSingleQuote bool, isDoubleQuote bool) bool {
	return char == ' ' && !isSingleQuote && !isDoubleQuote
}

func isSingleQuoteToggle(char byte, isDoubleQuote bool) bool {
	return char == '\'' && !isDoubleQuote
}

func isDoubleQuoteToggle(char byte, isSingleQuote bool) bool {
	return char == '"' && !isSingleQuote
}

func buildCommandFromTokens(tokens []string) Command {
	cmd := Command{
		Cmd: tokens[0],
	}
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