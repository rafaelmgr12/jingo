package jsongoparser

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input        string
		expectedType Value
	}{
		{
			input:        `{"key": "value"}`,
			expectedType: &Object{},
		},
		{
			input:        `{ "key1": true, "key2": false, "key3": null, "key4": "value", "key5": 101}`,
			expectedType: &Object{},
		},
		{
			input:        `{ "key": "value", "key-n": 101, "key-o": {}, "key-l": []}`,
			expectedType: &Object{},
		},
	}

	for i, tt := range tests {
		l := NewLexer(tt.input)
		p := NewParser(l)

		value, err := p.ParseJSON()
		if err != nil {
			t.Fatalf("Test %d: error parsing JSON: %v", i, err)
		}

		if value == nil {
			t.Fatalf("Test %d: parsed value is nil", i)
		}

		if _, ok := value.(*Object); !ok {
			t.Fatalf("Test %d: expected value type %T, got %T", i, tt.expectedType, value)
		}
	}
}

func TestLexerTokenization(t *testing.T) {
	tests := []struct {
		input    string
		expected []TokenType
	}{
		{
			input: `{"key": "value"}`,
			expected: []TokenType{
				TokenBraceOpen,
				TokenString,
				TokenColon,
				TokenString,
				TokenBraceClose,
				TokenEOF,
			},
		},
		{
			input: `[true, false, null, "string", 123]`,
			expected: []TokenType{
				TokenBracketOpen,
				TokenTrue,
				TokenComma,
				TokenFalse,
				TokenComma,
				TokenNull,
				TokenComma,
				TokenString,
				TokenComma,
				TokenNumber,
				TokenBracketClose,
				TokenEOF,
			},
		},
		{
			input: `{"key1": 100, "key2": 1.23, "key3": 2e10}`,
			expected: []TokenType{
				TokenBraceOpen,
				TokenString,
				TokenColon,
				TokenNumber,
				TokenComma,
				TokenString,
				TokenColon,
				TokenNumber,
				TokenComma,
				TokenString,
				TokenColon,
				TokenNumber,
				TokenBraceClose,
				TokenEOF,
			},
		},
	}

	for i, tt := range tests {
		l := NewLexer(tt.input)
		for _, expectedType := range tt.expected {
			token := l.NextToken()
			if token.Type != expectedType {
				t.Fatalf("Test %d: expected token type %q, got %q", i, expectedType, token.Type)
			}
		}
	}
}

func TestParserErrors(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr string
	}{
		{
			input:       `{"key": value}`,
			expectedErr: "expected string key",
		},
		{
			input:       `{"key" value}`,
			expectedErr: "expected :, got ILLEGAL",
		},
		{
			input:       `{"key": "value"`,
			expectedErr: "expected }, got EOF",
		},
		{
			input:       `{"key": "value",}`,
			expectedErr: "unexpected token ,",
		},
	}

	for i, tt := range tests {
		l := NewLexer(tt.input)
		p := NewParser(l)
		_, err := p.ParseJSON()
		errors := p.Errors()

		if err == nil {
			t.Errorf("Test %d: expected error but got none", i)
			continue
		}

		if !hasMatchingError(errors, tt.expectedErr) {
			t.Errorf("Test %d: expected error containing %q, got %v",
				i, tt.expectedErr, errors)
		}
	}
}

func TestComplexJSON(t *testing.T) {
	input := `{
		"key1": {
			"nestedKey1": "nestedValue1",
			"nestedKey2": [1, 2, {"deeplyNestedKey": "deeplyNestedValue"}]
		},
		"key2": "value2",
		"key3": 123.45,
		"key4": true
	}`

	l := NewLexer(input)
	p := NewParser(l)
	value, err := p.ParseJSON()
	if err != nil {
		t.Fatalf("Error parsing JSON: %v", err)
	}

	if obj, ok := value.(*Object); !ok || len(obj.Pairs) != 4 {
		t.Fatalf("Parsing resulted in wrong object structure: %+v", value)
	}
}

// hasMatchingError checks if any error in the list matches the expected error
func hasMatchingError(errors []string, expectedErr string) bool {
	for _, err := range errors {
		// Normalize both strings by trimming spaces and converting to lowercase
		normalizedErr := strings.ToLower(strings.TrimSpace(err))
		normalizedExpected := strings.ToLower(strings.TrimSpace(expectedErr))

		// Check if the normalized error contains the expected error string
		if strings.Contains(normalizedErr, normalizedExpected) {
			return true
		}
	}
	return false
}
