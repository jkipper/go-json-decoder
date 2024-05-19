package myjson_test

import (
	"testing"
	myjson "user/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInput(t *testing.T) {
	input := `
	{ "key": "value",
	"second_key": 10, 
	"nested": {"key":
	"value",
		
		"slice": [10, 20], 
	"bool": true}}
	`
	out, err := myjson.Decode(input)
	require.Nil(t, err)
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
	out, err := myjson.Decode(input)
	require.Nil(t, err)
	expected := map[string]interface{}{
		"key": "valu\\\"e",
	}
	assert.Equal(t, expected, out)
}

func TestMissingSeparator(t *testing.T) {
	input := `
	{"key1": "value1" "key2": "value2" }
	`
	_, err := myjson.Decode(input)
	expectedErr := &myjson.DecodeError{}
	require.ErrorAs(t, err, &expectedErr)
}

func TestObjInsideSlice(t *testing.T) {
	input := `
		{"key": [ {"key": "value"} ]}
	`
	out, err := myjson.Decode(input)
	require.Nil(t, err)

	expectedOutput := map[string]interface{}{
		"key": []interface{}{map[string]interface{}{"key": "value"}},
	}
	assert.Equal(t, expectedOutput, out)
}

func TestJustSlice(t *testing.T) {
	input := `
	[10,20,30,"string"]
	`
	expectedOutput := []interface{}{
		10, 20, 30, "string",
	}

	out, err := myjson.Decode(input)
	require.Nil(t, err)
	assert.Equal(t, expectedOutput, out)

}
