package myjson

import (
	"log"
	"strconv"
	"strings"
)

func parseObject(input []string) (map[string]interface{}, []string) {
	out := make(map[string]interface{})
	if input[0] == "{" {
		input = input[1:]
	}
	for input[0] != "}" {
		if len(out) > 0 {
			if input[0] == "," {
				input = input[1:]
			} else {
				log.Panicf("Missing seperator. Remaining texts: %s", input)
			}
		}

		key, remaining := parseKey(input)
		if remaining[0] != ":" {
			log.Panicf("Invalid token, expected seperator ':', got %s", remaining[0])
		}
		value, remaining := parseValue(remaining[1:])
		input = remaining
		out[key] = value
	}
	return out, input[1:]
}

func parseKey(input []string) (string, []string) {
	token := input[0]
	if token[0] == '"' {

		return token[1 : len(token)-1], input[1:]
	}
	return "", input
}

func parseValue(input []string) (interface{}, []string) {
	token := input[0]
	if token[0] == '"' {
		return token[1 : len(token)-1], input[1:]
	} else if token[0] == '{' {
		return parseObject(input[1:])
	} else if token[0] <= '9' && token[0] >= '0' {
		return parseNumber(token), input[1:]
	} else if token[0] == '[' {
		return parseSlice(input[1:])
	} else if token[0] == 't' || token[0] == 'f' {
		return parseBool(token), input[1:]
	} else {
		log.Panicf("Received invalid token %s, input: %v", token, input)
		return nil, nil
	}
}

func parseSlice(input []string) (interface{}, []string) {
	out := []interface{}{}
	for input[0] != "]" {
		if len(out) > 0 {
			if input[0] == "," {
				input = input[1:]
			} else {
				log.Panicf("Missing seperator between values, got %s", input)
			}
		}
		parsed, remaining := parseValue(input)
		if parsed == nil {
			panic("Received invalid token")
		}
		out = append(out, parsed)
		input = remaining
	}
	return out, input[1:]
}

func parseBool(token string) bool {
	if token == "true" {
		return true
	} else if token == "false" {
		return false
	} else {
		log.Panicf("Invalid value. Expected bool, got %s", token)
		return false
	}

}

func parseNumber(token string) interface{} {
	if strings.Contains(token, ".") {
		f, err := strconv.ParseFloat(token, 64)
		if err != nil {
			return nil
		} else {
			return f
		}
	} else {
		i, err := strconv.Atoi(token)
		if err != nil {
			return nil
		} else {
			return i
		}
	}
}
