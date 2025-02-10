package encoding_test

import (
	"reflect"
	"strings"
	"sync"
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
		_, err := encoding.Marshal(tt.input)
		if err != nil {
			t.Fatalf("Test %d (%s): error marshaling JSON: %v", i, tt.name, err)
		}

		var expectedMap, resultMap map[string]interface{}
		if !reflect.DeepEqual(expectedMap, resultMap) {
			t.Fatalf("Test %d (%s): expected %v, got %v", i, tt.name, expectedMap, resultMap)
		}
	}
}
func TestUnmarshalWithSizeLimit(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		maxSize     int
		shouldError bool
		errorCode   encoding.ErrorCode
		errorMsg    string
	}{
		{
			name:        "Within size limit",
			input:       []byte(`{"key": "value"}`),
			maxSize:     2048,
			shouldError: false,
		},
		{
			name:        "Exactly at size limit",
			input:       []byte(`{"key": "value"}`),
			maxSize:     1024,
			shouldError: false,
		},
		{
			name:        "Exceeds size limit",
			input:       []byte(`{"key": "value"}`),
			maxSize:     1073741825,
			shouldError: true,
			errorCode:   encoding.ErrInvalidOptions,
		},
		{
			name:        "Zero size limit",
			input:       []byte(`{"key": "value"}`),
			maxSize:     0,
			shouldError: false, // Should use default limit
		},
		{
			name:        "Negative size limit",
			input:       []byte(`{"key": "value"}`),
			maxSize:     -1,
			errorCode:   encoding.ErrInvalidOptions,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result map[string]interface{}

			var err error

			if tt.maxSize != 0 {
				err = encoding.Unmarshal(tt.input, &result, encoding.WithMaxSize(tt.maxSize))
			} else {
				err = encoding.Unmarshal(tt.input, &result)
			}

			if tt.shouldError {
				if err == nil {
					t.Errorf("%s: Expected error but got none", tt.name)
				} else {
					checkJSONError(t, err, tt.errorCode, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("%s: Unexpected error: %v", tt.name, err)
				}
			}
		})
	}
}

func TestMarshalWithSizeLimit(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		maxSize     int
		shouldError bool
		errorMsg    string
	}{
		{
			name: "Within size limit",
			input: map[string]string{
				"key": "value",
			},
			maxSize:     1024,
			shouldError: false,
		},
		{
			name: "Exceeds size limit",
			input: map[string]string{
				"key": "value",
			},
			maxSize:     1073741825,
			shouldError: true,
			errorMsg:    "max size 1073741825 exceeds maximum allowed size 1073741824",
		},
		{
			name: "Large nested structure",
			input: map[string]interface{}{
				"array": make([]int, 1000),
				"nested": map[string]string{
					"key": strings.Repeat("long string", 100),
				},
			},
			maxSize:     1024,
			shouldError: true,
			errorMsg:    "size_exceeded: size 3131 exceeds limit 1024",
		},
		{
			name: "Default size limit",
			input: map[string]interface{}{
				"key": strings.Repeat("a", 1000),
			},
			maxSize:     0, // Use default
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			var result []byte

			if tt.maxSize != 0 {
				result, err = encoding.Marshal(tt.input, encoding.WithMaxSize(tt.maxSize))
			} else {
				result, err = encoding.Marshal(tt.input)
			}

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if result == nil {
					t.Error("Expected result but got nil")
				}
			}
		})
	}
}

func TestOptionsConfiguration(t *testing.T) {
	tests := []struct {
		name          string
		options       []encoding.Option
		expectedSize  int
		input         []byte
		shouldSucceed bool
	}{
		{
			name:          "Default options",
			options:       nil,
			expectedSize:  encoding.DefaultMaxSize,
			input:         []byte(`{"key": "value"}`),
			shouldSucceed: true,
		},
		{
			name: "Custom size",
			options: []encoding.Option{
				encoding.WithMaxSize(2000),
			},
			expectedSize:  2000,
			input:         []byte(`{"key": "value"}`),
			shouldSucceed: true,
		},
		{
			name: "Multiple options (last wins)",
			options: []encoding.Option{
				encoding.WithMaxSize(1024),
				encoding.WithMaxSize(2048),
			},
			expectedSize:  2048,
			input:         []byte(`{"key": "value"}`),
			shouldSucceed: true,
		},
		{
			name: "Invalid size (negative)",
			options: []encoding.Option{
				encoding.WithMaxSize(-1),
			},
			expectedSize:  encoding.DefaultMaxSize,
			input:         []byte(`{"key": "value"}`),
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result map[string]interface{}
			err := encoding.Unmarshal(tt.input, &result, tt.options...)

			if tt.shouldSucceed {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Error("Expected error but got none")
				}
			}
		})
	}
}

func TestConcurrentOptionUsage(t *testing.T) {
	input := []byte(`{"key": "value"}`)

	var wg sync.WaitGroup

	iterations := 100
	minimumSize := 1024

	wg.Add(iterations)

	for i := minimumSize; i < minimumSize+iterations; i++ {
		go func(size int) {
			defer wg.Done()

			var result map[string]interface{}

			err := encoding.Unmarshal(input, &result, encoding.WithMaxSize(size))
			if size >= len(input) && err != nil {
				t.Errorf("Unexpected error for size %d: %v", size, err)
			}

			if size < len(input) && err == nil {
				t.Errorf("Expected error for size %d but got none", size)
			}
		}(i + 10)
	}

	wg.Wait()
}

func checkJSONError(t *testing.T, err error, expectedCode encoding.ErrorCode, expectedMsg string) {
	t.Helper()

	if jsonErr, ok := err.(*encoding.JSONError); ok {
		if jsonErr.Code != expectedCode {
			t.Errorf("expected error code %q, got %q", expectedCode, jsonErr.Code)
		}

		if expectedMsg != "" && !strings.Contains(jsonErr.Error(), expectedMsg) {
			t.Errorf("expected error message to contain %q, got %q", expectedMsg, jsonErr.Error())
		}
	} else {
		t.Errorf("expected JSONError, got %T: %v", err, err)
	}
}
