package myjson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLex(t *testing.T) {
	input := `
	{ "key": "value", "second_key": 10}
	`
	out := lex(input)
	expected := []string{"{", "\"key\"", ":", "\"value\"", ",", "\"second_key\"", ":", "10", "}"}
	assert.Equal(t, expected, out)
}

func TestEscapedString(t *testing.T) {
	input := `
	{"key": "va\"lue" }
	`
	out := lex(input)
	expected := []string{"{", "\"key\"", ":", "\"va\\\"lue\"", "}"}
	assert.Equal(t, expected, out)
}
