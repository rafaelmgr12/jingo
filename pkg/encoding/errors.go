package encoding

import "fmt"

// ErrorCode represents specific error types that can occur during encoding
type ErrorCode string

const (
	// Size-related errors
	ErrSizeExceeded ErrorCode = "size_exceeded"

	// Parse-related errors
	ErrInvalidJSON    ErrorCode = "invalid_json"
	ErrUnexpectedType ErrorCode = "unexpected_type"
	ErrInvalidValue   ErrorCode = "invalid_value"

	// Marshal-related errors
	ErrMarshalFailure  ErrorCode = "marshal_failure"
	ErrUnsupportedType ErrorCode = "unsupported_type"

	// Unmarshal-related errors
	ErrUnmarshalFailure ErrorCode = "unmarshal_failure"
	ErrInvalidTarget    ErrorCode = "invalid_target"

	// Configuration errors
	ErrInvalidOptions ErrorCode = "invalid_options"
)

// JSONError represents a structured error that occurs during JSON processing
type JSONError struct {
	// Code identifies the specific type of error
	Code ErrorCode

	// Message provides a human-readable description of the error
	Message string

	// Path represents the JSON path where the error occurred (if applicable)
	Path string

	// Value contains the problematic value (if applicable)
	Value interface{}

	// Cause is the underlying error that caused this error (if any)
	Cause error
}

// Error implements the error interface with a formatted message
func (e *JSONError) Error() string {
	msg := string(e.Code)

	if e.Message != "" {
		msg += ": " + e.Message
	}

	if e.Path != "" {
		msg += fmt.Sprintf(" (at %s)", e.Path)
	}

	if e.Cause != nil {
		msg += fmt.Sprintf(" (caused by: %v)", e.Cause)
	}

	return msg
}

// Unwrap implements the unwrap interface for error chains
func (e *JSONError) Unwrap() error {
	return e.Cause
}

// NewJSONError creates a new JSONError with the given code and message
func NewJSONError(code ErrorCode, msg string) *JSONError {
	return &JSONError{
		Code:    code,
		Message: msg,
	}
}

// WithPath adds a JSON path to the error
func (e *JSONError) WithPath(path string) *JSONError {
	e.Path = path

	return e
}

// WithValue adds a problematic value to the error
func (e *JSONError) WithValue(value interface{}) *JSONError {
	e.Value = value

	return e
}

// WithCause adds an underlying cause to the error
func (e *JSONError) WithCause(err error) *JSONError {
	e.Cause = err

	return e
}

// Error creation helper functions
func NewSizeExceededError(size, limit int) *JSONError {
	return NewJSONError(ErrSizeExceeded,
		fmt.Sprintf("size %d exceeds limit %d", size, limit))
}

func NewInvalidTargetError(typ string) *JSONError {
	return NewJSONError(ErrInvalidTarget,
		fmt.Sprintf("invalid target type: %s", typ))
}

func NewUnsupportedTypeError(typ string) *JSONError {
	return NewJSONError(ErrUnsupportedType,
		fmt.Sprintf("unsupported type: %s", typ))
}

func NewUnmarshalTypeError(expected, got string) *JSONError {
	return NewJSONError(ErrUnmarshalFailure,
		fmt.Sprintf("cannot unmarshal %s into %s", got, expected))
}
