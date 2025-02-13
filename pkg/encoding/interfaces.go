package encoding

import "io"

// JSONDecoder defines the interface for decoding JSON values from a stream
type JSONDecoder interface {
	// Decode reads the next JSON-encoded value from its input and stores it in v
	Decode(v interface{}) error
	// More reports whether there is another value in the input stream
	More() bool
	// BufferSize returns the size of the underlying buffer
	BufferSize() int
}

// JSONEncoder defines the interface for encoding JSON values to a stream
type JSONEncoder interface {
	// Encode writes the JSON encoding of v to the stream
	Encode(v interface{}) error
	// SetIndent sets the indentation string for pretty-printing
	SetIndent(prefix, indent string)
	// Flush ensures all buffered data is written to the underlying writer
	Flush() error
}

// JSONStreamProcessor combines encoding and decoding capabilities
type JSONStreamProcessor interface {
	JSONEncoder
	JSONDecoder
}

// Writer extends io.Writer with JSON-specific writing capabilities
type Writer interface {
	io.Writer
	// WriteJSON writes a JSON-encoded value
	WriteJSON(v interface{}) error
	// FormatJSON controls the output formatting
	FormatJSON(pretty bool)
}

// Reader extends io.Reader with JSON-specific reading capabilities
type Reader interface {
	io.Reader
	// ReadJSON reads a JSON-encoded value into v
	ReadJSON(v interface{}) error
	// Skip skips the next JSON value in the stream
	Skip() error
}

// Stream defines the interface for a bi-directional JSON stream
type Stream interface {
	Writer
	Reader
	// Close closes the stream and releases any resources
	Close() error
}
