package myjson

import "log"

func Decode(input string) (map[string]interface{}, error) {
	tokens := lex(input)
	data, remainder, err := parseObject(tokens)
	if err != nil {
		return nil, err
	}
	if len(remainder) > 0 {
		log.Panicf("Failed to parse complete document. Remainging tokens: %v", remainder)
	}
	return data, nil

}
