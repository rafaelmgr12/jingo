package encoding_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/rafaelmgr12/jingo/pkg/encoding"
)

func TestUnmarshalWithByteInput(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected interface{}
	}{
		{
			name:  "Simple JSON",
			input: []byte(`{"key": "value"}`),
			expected: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:  "Complex JSON",
			input: []byte(`{"key1": true, "key2": false, "key3": null, "key4": "value", "key5": 101}`),
			expected: map[string]interface{}{
				"key1": true,
				"key2": false,
				"key3": nil,
				"key4": "value",
				"key5": int64(101),
			},
		},
		{
			name:  "Complex JSON with nested objects",
			input: []byte(`{"key": "value", "key-n": 101, "key-o": {}, "key-l": []}`),
			expected: map[string]interface{}{
				"key":   "value",
				"key-n": int64(101),
				"key-o": map[string]interface{}{},
				"key-l": []interface{}{},
			},
		},
		{
			name:  "Array of objects",
			input: []byte(`[{"key1": "value1"}, {"key2": "value2"}]`),
			expected: []interface{}{
				map[string]interface{}{"key1": "value1"},
				map[string]interface{}{"key2": "value2"},
			},
		},
		{
			name:  "Nested arrays",
			input: []byte(`{"key": [[1, 2], [3, 4]]}`),
			expected: map[string]interface{}{
				"key": []interface{}{
					[]interface{}{int64(1), int64(2)},
					[]interface{}{int64(3), int64(4)},
				},
			},
		},
		{
			name:     "Empty JSON object",
			input:    []byte(`{}`),
			expected: map[string]interface{}{},
		},
		{
			name:     "Empty JSON array",
			input:    []byte(`[]`),
			expected: []interface{}{},
		},
		{
			name:  "JSON with special characters",
			input: []byte(`{"key": "value with special characters !@#$%^&*()"}`),
			expected: map[string]interface{}{
				"key": "value with special characters !@#$%^&*()",
			},
		},
		{
			name:  "JSON with unicode characters",
			input: []byte(`{"key": "こんにちは"}`),
			expected: map[string]interface{}{
				"key": "こんにちは",
			},
		},
	}

	for i, tt := range tests {
		var result interface{}
		if tt.name == "Array of objects" || tt.name == "Empty JSON array" {
			result = []interface{}{}
		} else {
			result = map[string]interface{}{}
		}

		err := encoding.Unmarshal(tt.input, &result)
		if err != nil {
			t.Fatalf("Test %d (%s): error unmarshaling JSON: %v", i, tt.name, err)
		}

		if !reflect.DeepEqual(tt.expected, result) {
			// iterate trhour ther result interfae and print the tyoes
			for k, v := range result.(map[string]interface{}) {
				t.Logf("Key: %s, Value: %v, Type: %T", k, v, v)
			}
			// itereater over the tt.expected interface and print the types
			for k, v := range tt.expected.(map[string]interface{}) {
				t.Logf("Key: %s, Value: %v, Type: %T", k, v, v)
			}
			t.Fatalf("Test %d (%s): expected %v, got %v", i, tt.name, tt.expected, result)
		}

	}
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "Simple map",
			input: map[string]interface{}{
				"key": "value",
			},
			expected: `{"key":"value"}`,
		},
		{
			name: "Complex map",
			input: map[string]interface{}{
				"key1": true,
				"key2": false,
				"key3": nil,
				"key4": "value",
				"key5": 101,
			},
			expected: `{"key1":true,"key2":false,"key3":null,"key4":"value","key5":101}`,
		},
		{
			name: "Nested map",
			input: map[string]interface{}{
				"key": map[string]interface{}{
					"nestedKey": "nestedValue",
				},
			},
			expected: `{"key":{"nestedKey":"nestedValue"}}`,
		},
		{
			name: "Array of objects",
			input: []interface{}{
				map[string]interface{}{"key1": "value1"},
				map[string]interface{}{"key2": "value2"},
			},
			expected: `[{"key1":"value1"},{"key2":"value2"}]`,
		},
		{
			name: "Nested arrays",
			input: map[string]interface{}{
				"key": []interface{}{
					[]interface{}{1, 2},
					[]interface{}{3, 4},
				},
			},
			expected: `{"key":[[1,2],[3,4]]}`,
		},
		{
			name:     "Empty map",
			input:    map[string]interface{}{},
			expected: `{}`,
		},
		{
			name:     "Empty array",
			input:    []interface{}{},
			expected: `[]`,
		},
		{
			name: "Special characters",
			input: map[string]interface{}{
				"key": "value with special characters !@#$%^&*()",
			},
			expected: `{"key":"value with special characters !@#$%^&*()"}`,
		},
		{
			name: "Unicode characters",
			input: map[string]interface{}{
				"key": "こんにちは",
			},
			expected: `{"key":"こんにちは"}`,
		},
	}

	for i, tt := range tests {
		result, err := encoding.Marshal(tt.input)
		if err != nil {
			t.Fatalf("Test %d (%s): error marshaling JSON: %v", i, tt.name, err)
		}

		var expectedMap, resultMap map[string]interface{}
		json.Unmarshal([]byte(tt.expected), &expectedMap)
		json.Unmarshal(result, &resultMap)

		if !reflect.DeepEqual(expectedMap, resultMap) {
			t.Fatalf("Test %d (%s): expected %v, got %v", i, tt.name, expectedMap, resultMap)
		}
	}
}
