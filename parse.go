package myjson

import (
	"fmt"
	"strconv"
	"strings"
)

type DecodeError struct {
	Expected string
	Got      string
}

func (err *DecodeError) Error() string {
	return "Expected '" + err.Expected + "' got '" + err.Got + "'"
}

func parseObject(input []string) (interface{}, []string, error) {
	out := make(map[string]interface{})
	if input[0] == "{" {
		input = input[1:]
	} else {
		return nil, input, nil
	}
	for input[0] != "}" {
		if len(out) > 0 {
			if input[0] == "," {
				input = input[1:]
			} else {
				return nil, nil, fmt.Errorf("Missing seperator: %w", &DecodeError{Expected: ",", Got: input[0]})
			}
		}

		key, remaining, err := parseKey(input)
		if err != nil {
			return nil, nil, err
		}
		if remaining[0] != ":" {
			return nil, nil, fmt.Errorf("Missing seperator between key and value: %w", &DecodeError{Expected: ":", Got: remaining[0]})
		}
		value, remaining, err := parseValue(remaining[1:])
		if err != nil {
			return nil, nil, err
		}
		input = remaining
		out[key] = value
	}
	return out, input[1:], nil
}

func parseKey(input []string) (string, []string, error) {
	got, remaining, err := parseString(input)
	if err != nil {
		return "", nil, err
	}
	if got == "" {
		return "", input, &DecodeError{Expected: "String", Got: input[0]}
	}

	return got.(string), remaining, nil
}

func parseString(input []string) (interface{}, []string, error) {
	token := input[0]
	if token[0] == '"' {
		return token[1 : len(token)-1], input[1:], nil
	}
	return "", input, nil
}

func parseValue(input []string) (interface{}, []string, error) {
	parsers := []func([]string) (interface{}, []string, error){
		parseString, parseObject, parseNumber, parseSlice, parseBool,
	}
	for _, parser := range parsers {
		token, remaining, err := parser(input)
		if err != nil {
			return nil, input, err
		}
		if len(remaining) < len(input) {
			return token, remaining, nil
		}
	}
	err := fmt.Errorf("Received invalid token %v", input[0])
	return nil, nil, err
}

func parseSlice(input []string) (interface{}, []string, error) {
	if input[0] != "[" {
		return nil, input, nil
	} else {
		input = input[1:]
	}

	out := []interface{}{}
	for input[0] != "]" {
		if len(out) > 0 {
			if input[0] == "," {
				input = input[1:]
			} else {
				err := fmt.Errorf("Missing separator: %w", &DecodeError{Expected: ",", Got: input[0]})
				return nil, nil, err
			}
		}
		parsed, remaining, err := parseValue(input)
		if err != nil {
			return nil, nil, err
		}
		out = append(out, parsed)
		input = remaining
	}
	return out, input[1:], nil
}

func parseBool(input []string) (interface{}, []string, error) {
	token := input[0]
	if token == "true" {
		return true, input[1:], nil
	} else if token == "false" {
		return false, input[1:], nil
	} else {
		return nil, input, nil
	}

}

func parseNumber(input []string) (interface{}, []string, error) {
	token := input[0]
	isNumber := token[0] >= '0' && token[0] <= '9'
	if !isNumber {
		return nil, input, nil
	}
	if strings.Contains(token, ".") {
		f, err := strconv.ParseFloat(token, 64)
		if err != nil {
			return nil, input, err
		} else {
			return f, input[1:], nil
		}
	} else {
		i, err := strconv.Atoi(token)
		if err != nil {
			return nil, input, err
		} else {
			return i, input[1:], nil
		}
	}
}
