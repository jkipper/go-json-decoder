package myjson_test

import (
	"testing"
	myjson "user/json"

	"github.com/stretchr/testify/assert"
)

func TestInput(t *testing.T) {
	input := `
	{ "key": "value", "second_key": 10, "nested": {"key": "value", "slice": [10, 20], "bool": true}}
	`
	out := myjson.Decode(input)
	expected := map[string]interface{}{
		"key":        "value",
		"second_key": 10,
		"nested": map[string]interface{}{
			"key":   "value",
			"slice": []interface{}{10, 20},
			"bool":  true,
		},
	}
	assert.Equal(t, expected, out)
}

func TestEscapedString(t *testing.T) {
	input := `
	{
		"key" : "valu\"e"
	}
	`
	out := myjson.Decode(input)
	expected := map[string]interface{}{
		"key": "valu\\\"e",
	}
	assert.Equal(t, expected, out)
}
