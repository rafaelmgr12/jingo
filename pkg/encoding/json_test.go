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

func TestUnmarshalWithSizeLimit(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		maxSize     int
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Within size limit",
			input:       []byte(`{"key": "value"}`),
			maxSize:     100,
			shouldError: false,
		},
		{
			name:        "Exactly at size limit",
			input:       []byte(`{"key": "value"}`),
			maxSize:     16, // 15 bytes in the previous assumption overlooked newline
			shouldError: false,
		},
		{
			name:        "Exceeds size limit",
			input:       []byte(`{"key": "value"}`),
			maxSize:     10,
			shouldError: true,
			errorMsg:    "input JSON size (16 bytes) exceeds maximum allowed size (10 bytes)", // Fixing after test observations
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
			shouldError: false, // Should use default limit
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
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("%s: Expected error message %q, got %q", tt.name, tt.errorMsg, err.Error())
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
			maxSize:     100,
			shouldError: false,
		},
		{
			name: "Exceeds size limit",
			input: map[string]string{
				"key": "value",
			},
			maxSize:     10,
			shouldError: true,
			errorMsg:    "marshaled JSON size",
		},
		{
			name: "Large nested structure",
			input: map[string]interface{}{
				"array": make([]int, 1000),
				"nested": map[string]string{
					"key": strings.Repeat("long string", 100),
				},
			},
			maxSize:     1000,
			shouldError: true,
			errorMsg:    "exceeds maximum allowed size",
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
				encoding.WithMaxSize(100),
			},
			expectedSize:  100,
			input:         []byte(`{"key": "value"}`),
			shouldSucceed: true,
		},
		{
			name: "Multiple options (last wins)",
			options: []encoding.Option{
				encoding.WithMaxSize(100),
				encoding.WithMaxSize(200),
			},
			expectedSize:  200,
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
			shouldSucceed: true,
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

	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
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
