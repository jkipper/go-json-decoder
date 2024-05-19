package myjson

import (
	"fmt"
	"log"
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

var _ error = &DecodeError{}

func parseObject(input []string) (map[string]interface{}, []string, error) {
	out := make(map[string]interface{})
	if input[0] == "{" {
		input = input[1:]
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
			return nil, nil, fmt.Errorf("Invalid Token: %w", &DecodeError{Expected: ":", Got: remaining[0]})
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
	token := input[0]
	if token[0] == '"' {
		return token[1 : len(token)-1], input[1:], nil
	}
	return "", input, &DecodeError{Expected: "String", Got: token}
}

func parseValue(input []string) (interface{}, []string, error) {
	token := input[0]
	if token[0] == '"' {
		return token[1 : len(token)-1], input[1:], nil
	} else if token[0] == '{' {
		return parseObject(input[1:])
	} else if token[0] <= '9' && token[0] >= '0' {
		number, err := parseNumber(token)
		if err != nil {
			return nil, nil, err
		}
		return number, input[1:], nil
	} else if token[0] == '[' {
		return parseSlice(input[1:])
	} else if token[0] == 't' || token[0] == 'f' {
		return parseBool(token), input[1:], nil
	} else {
		err := fmt.Errorf("Received invalid token %v", token)
		return nil, nil, err
	}
}

func parseSlice(input []string) (interface{}, []string, error) {
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

func parseBool(token string) bool {
	if token == "true" {
		return true
	} else if token == "false" {
		return false
	} else {
		// should never happen
		log.Panicf("Invalid value. Expected bool, got %s", token)
		return false
	}

}

func parseNumber(token string) (interface{}, error) {
	if strings.Contains(token, ".") {
		f, err := strconv.ParseFloat(token, 64)
		if err != nil {
			return nil, err
		} else {
			return f, nil
		}
	} else {
		i, err := strconv.Atoi(token)
		if err != nil {
			return nil, err
		} else {
			return i, nil
		}
	}
}
