package myjson

import (
	"log"
	"slices"
	"strings"
)

func lex(input string) []string {
	tokens := []string{}
	lexers := []func(string) (string, string){
		lexKeyToken, lexStr, lexNumber, lexBool, lexSkip,
	}

	for len(input) > 0 {
		for _, lexer := range lexers {
			token, remaining := lexer(input)
			if len(remaining) < len(input) {
				if len(token) > 0 {
					tokens = append(tokens, token)
				}
				input = remaining
				break
			}
		}
	}
	return tokens
}

func lexSkip(input string) (string, string) {
	skippable := []rune{' ', '\n', '\t', '\r'}
	if slices.Contains(skippable, rune(input[0])) {
		return "", input[1:]
	} else {
		log.Panicf("Encountered unkown character %s", input)
		return "", input
	}
}

func lexKeyToken(input string) (string, string) {
	maybeToken := input[0]
	keyTokens := []rune{'{', '}', ',', ':', '[', ']'}
	if slices.Contains(keyTokens, rune(maybeToken)) {
		return string(maybeToken), input[1:]
	}
	return "", input
}

func lexStr(input string) (string, string) {
	if input[0] != '"' {
		return "", input
	}
	output := ""
	for idx, v := range input {
		output += string(v)
		if v == '"' && idx > 0 && input[idx-1] != '\\' {
			return output, input[len(output):]
		}
	}
	return output, input[len(output):]
}

func lexNumber(input string) (string, string) {
	out := ""
	for _, v := range input {
		if v >= '0' && v <= '9' {
			out += string(v)
		} else {
			break
		}
	}
	return out, input[len(out):]
}

func lexBool(input string) (string, string) {
	if strings.HasPrefix(input, "true") {
		return "true", input[len("true"):]
	} else if strings.HasPrefix(input, "false") {
		return "false", input[len("false"):]
	} else {
		return "", input
	}
}
