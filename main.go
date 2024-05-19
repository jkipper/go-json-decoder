package myjson

import (
	"fmt"
)

func Decode(input string) (interface{}, error) {
	tokens := lex(input)
	data, remainder, err := Parse(tokens)
	if err != nil {
		return nil, err
	}
	if len(remainder) > 0 {
		return nil, fmt.Errorf("Failed to parse complete document. Remaining tokens: %v", remainder)
	}
	return data, nil
}
