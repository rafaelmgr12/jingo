package encoding_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/rafaelmgr12/jingo/pkg/encoding"
)

func TestNewDecoder(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		options     []encoding.Option
		expectedErr string
	}{
		{
			name:    "Default options",
			input:   `{"key": "value"}`,
			options: nil,
		},
		{
			name:        "Invalid buffer size",
			input:       `{"key": "value"}`,
			options:     []encoding.Option{encoding.WithBufferSize(-1)},
			expectedErr: "invalid option",
		},
		{
			name:    "Custom buffer size",
			input:   `{"key": "value"}`,
			options: []encoding.Option{encoding.WithBufferSize(8192)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reader := strings.NewReader(tt.input)
			decoder, err := encoding.NewDecoder(reader, tt.options...)

			if tt.expectedErr != "" {
				if err == nil {
					t.Fatalf("Expected error %v, but got none", tt.expectedErr)
				}

				if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Fatalf("Expected error to contain %v, but got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				var result map[string]interface{}

				if err := decoder.Decode(&result); err != nil {
					t.Fatalf("Failed to decode: %v", err)
				}

				if result["key"] != "value" {
					t.Fatalf("Expected key to be 'value', got %v", result["key"])
				}
			}
		})
	}
}

func TestNewEncoder(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		options     []encoding.Option
		expectedErr string
	}{
		{
			name:    "Default options",
			input:   map[string]string{"key": "value"},
			options: nil,
		},
		{
			name:        "Invalid buffer size",
			input:       map[string]string{"key": "value"},
			options:     []encoding.Option{encoding.WithBufferSize(-1)},
			expectedErr: "invalid option",
		},
		{
			name:    "Custom buffer size",
			input:   map[string]string{"key": "value"},
			options: []encoding.Option{encoding.WithBufferSize(8192)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer

			encoder, err := encoding.NewEncoder(&buffer, tt.options...)

			if tt.expectedErr != "" {
				if err == nil {
					t.Fatalf("Expected error %v, but got none", tt.expectedErr)
				}

				if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Fatalf("Expected error to contain %v, but got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				if err := encoder.Encode(tt.input); err != nil {
					t.Fatalf("Failed to encode: %v", err)
				}

				var expectedMap, resultMap map[string]interface{}

				inputBytes, err := json.Marshal(tt.input)
				if err != nil {
					t.Fatalf("Failed to marshal input: %v", err)
				}

				if err := json.Unmarshal(inputBytes, &expectedMap); err != nil {
					t.Fatalf("Failed to unmarshal expected JSON: %v", err)
				}

				if err := json.Unmarshal(buffer.Bytes(), &resultMap); err != nil {
					t.Fatalf("Failed to unmarshal result JSON: %v", err)
				}

				if !reflect.DeepEqual(expectedMap, resultMap) {
					t.Fatalf("Expected %v, got %v", expectedMap, resultMap)
				}
			}
		})
	}
}

func TestOptions(t *testing.T) {
	tests := []struct {
		name        string
		options     []encoding.Option
		expectedErr string
	}{
		{
			name:    "Default options",
			options: nil,
		},
		{
			name:        "Max size too small",
			options:     []encoding.Option{encoding.WithMaxSize(512)},
			expectedErr: "max size 512 is below minimum allowed size 1024",
		},
		{
			name:        "Invalid buffer size",
			options:     []encoding.Option{encoding.WithBufferSize(-1)},
			expectedErr: "buffer size must be positive, got -1",
		},
		{
			name: "Valid buffer size",
			options: []encoding.Option{
				encoding.WithBufferSize(8192),
			},
		},
		{
			name: "Valid strict mode",
			options: []encoding.Option{
				encoding.WithStrictMode(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encoding.NewDecoder(strings.NewReader(`{"key": "value"}`), tt.options...)

			if tt.expectedErr != "" {
				if err == nil {
					t.Fatalf("Expected error %v, but got none", tt.expectedErr)
				}

				if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Fatalf("Expected error to contain %v, but got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
			}
		})
	}
}
