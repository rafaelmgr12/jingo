package jsongoparser

import (
	"fmt"
	"math"
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

func TestNumberParsing(t *testing.T) {
	tests := []struct {
		input      string
		wantInt    bool
		wantValue  interface{} // can be int64 or float64
		shouldFail bool
	}{
		// Valid cases
		{
			input:      `{"num": 123}`,
			wantInt:    true,
			wantValue:  int64(123),
			shouldFail: false,
		},
		{
			input:      `{"num": 123.456}`,
			wantInt:    false,
			wantValue:  float64(123.456),
			shouldFail: false,
		},
		{
			input:      `{"num": -123}`,
			wantInt:    true,
			wantValue:  int64(-123),
			shouldFail: false,
		},
		{
			input:      `{"num": 1e5}`,
			wantInt:    false,
			wantValue:  float64(100000),
			shouldFail: false,
		},
		{
			input:      `{"num": 1.2e-3}`,
			wantInt:    false,
			wantValue:  float64(0.0012),
			shouldFail: false,
		},
		// Invalid cases
		{
			input:      `{"num": 01234}`, // Leading zeros not allowed
			shouldFail: true,
		},
		{
			input:      `{"num": .123}`, // Must start with digit
			shouldFail: true,
		},
		{
			input:      `{"num": 123.}`, // Must have digits after decimal
			shouldFail: true,
		},
		{
			input:      `{"num": -}`, // Must have digits after minus
			shouldFail: true,
		},
		{
			input:      `{"num": 1.2e}`, // Must have digits after exponent
			shouldFail: true,
		},
		{
			input:      `{"num": 1.2e-}`, // Must have digits after exponent sign
			shouldFail: true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Case %d: %s", i, tt.input), func(t *testing.T) {
			l := NewLexer(tt.input)
			p := NewParser(l)
			value, err := p.ParseJSON()

			if tt.shouldFail {
				if err == nil {
					t.Errorf("Expected error for input %s but got none", tt.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			obj, ok := value.(*Object)
			if !ok {
				t.Fatalf("Expected Object, got %T", value)
			}

			num, ok := obj.Pairs["num"].(*NumberLiteral)
			if !ok {
				t.Fatalf("Expected NumberLiteral, got %T", obj.Pairs["num"])
			}

			if num.IsInt != tt.wantInt {
				t.Errorf("IsInt = %v, want %v", num.IsInt, tt.wantInt)
			}

			if tt.wantInt {
				if num.Int != tt.wantValue.(int64) {
					t.Errorf("Int = %d, want %d", num.Int, tt.wantValue.(int64))
				}
			} else {
				if math.Abs(num.Float-tt.wantValue.(float64)) > 1e-10 {
					t.Errorf("Float = %g, want %g", num.Float, tt.wantValue.(float64))
				}
			}
		})
	}
}

func FuzzParseJSON(f *testing.F) {
	// Add initial seed corpus
	f.Add(`{"key": "value"}`)
	f.Add(`[1, 2, 3]`)
	f.Add(`{"nested": {"key": "value"}, "array": [1, 2, 3]}`)
	f.Add(`{"number": 12345}`)
	f.Add(`{"boolean": true}`)
	f.Add(`{"nullValue": null}`)
	f.Add(`{"escapedString": "Line1\nLine2"}`)
	f.Add(`{"unicode": "こんにちは"}`)
	f.Add(`{"specialChars": "!@#$%^&*()_+-=~"}`)
	f.Add(`{"emptyObject": {}}`)
	f.Add(`{"emptyArray": []}`)
	f.Add(`{"mixedArray": [1, "two", true, null, {"key": "value"}]}`)
	f.Add(`{"deeplyNested": {"level1": {"level2": {"level3": {"level4": "value"}}}}}`)

	f.Fuzz(func(t *testing.T, input string) {
		lexer := NewLexer(input)
		parser := NewParser(lexer)
		parsed, err := parser.ParseJSON()

		if err != nil {
			// Check for specific error types or messages
			if !isExpectedError(err) {
				t.Errorf("Unexpected error parsing JSON: %v", err)
			}
		} else {
			// Validate the parsed output for known valid inputs
			if !isValidParsedOutput(parsed) {
				t.Errorf("Parsed output is invalid for input: %s", input)
			}
		}
	})
}

func BenchmarkParseJSON(b *testing.B) {
	input := `{
		"key1": "value1",
		"key2": 123,
		"key3": [1, 2, 3],
		"key4": {"nestedKey": "nestedValue"},
		"key5": true,
		"key6": null
	}`

	for i := 0; i < b.N; i++ {
		lexer := NewLexer(input)
		parser := NewParser(lexer)
		_, err := parser.ParseJSON()
		if err != nil {
			b.Fatalf("Error parsing JSON: %v", err)
		}
	}
}

// isExpectedError checks if the error is one of the expected errors
func isExpectedError(err error) bool {
	expectedErrors := []string{
		"unexpected token",
		"invalid character",
		"unterminated string",
		"invalid number format",
	}
	for _, expectedErr := range expectedErrors {
		if strings.Contains(err.Error(), expectedErr) {
			return true
		}
	}
	return false
}

// isValidParsedOutput validates the parsed output for known valid inputs
func isValidParsedOutput(parsed interface{}) bool {
	switch parsed.(type) {
	case *Object, *Array:
		return true
	}
	return false
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
