package myjson

import "log"

func Decode(input string) map[string]interface{} {
	tokens := lex(input)
	data, remainder := parseObject(tokens)
	if len(remainder) > 0 {
		log.Panicf("Failed to parse complete document. Remainging tokens: %v", remainder)
	}
	return data

}
