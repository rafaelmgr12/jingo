package encoding

import (
	"bufio"
	"io"
	"sync"
)

// streamEncoder provides a concrete implementation of JSONEncoder interface.
// It handles the encoding of JSON values to an output stream with support
// for pretty printing and proper error handling.
type streamEncoder struct {
	writer     *bufio.Writer
	options    *Options
	mutex      sync.Mutex
	prefix     string
	indent     string
	bufferSize int
}

// NewEncoder creates a new JSONEncoder implementation.
// It accepts an io.Writer and optional configuration options.
func NewEncoder(w io.Writer, opts ...Option) (JSONEncoder, error) {
	options, err := applyOptions(opts...)
	if err != nil {
		return nil, NewJSONError(ErrInvalidOptions, "invalid encoder options").WithCause(err)
	}

	bufferSize := 4096
	if options.BufferSize > 0 {
		bufferSize = options.BufferSize
	}

	return &streamEncoder{
		writer:     bufio.NewWriterSize(w, bufferSize),
		options:    options,
		bufferSize: bufferSize,
	}, nil
}

// Encode implements JSONEncoder.Encode.
// It writes the JSON encoding of v to the stream.
func (e *streamEncoder) Encode(v interface{}) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	var data []byte

	var err error

	data, err = Marshal(v)

	if err != nil {
		return NewJSONError(ErrMarshalFailure, "failed to marshal value for stream").
			WithCause(err).
			WithValue(v)
	}

	if !e.options.DisableSizeLimit && len(data) > e.options.MaxSize {
		return NewSizeExceededError(len(data), e.options.MaxSize)
	}

	if _, err := e.writer.Write(data); err != nil {
		return NewJSONError(ErrMarshalFailure, "failed to write to stream").WithCause(err)
	}

	if err := e.writer.WriteByte('\n'); err != nil {
		return NewJSONError(ErrMarshalFailure, "failed to write newline to stream").WithCause(err)
	}

	return e.Flush()
}

// SetIndent implements JSONEncoder.SetIndent.
// It configures the encoder's indentation settings for pretty printing.
func (e *streamEncoder) SetIndent(prefix, indent string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.prefix = prefix
	e.indent = indent
}

// Flush implements JSONEncoder.Flush.
// It ensures all buffered data is written to the underlying writer.
func (e *streamEncoder) Flush() error {
	if err := e.writer.Flush(); err != nil {
		return NewJSONError(ErrMarshalFailure, "failed to flush stream").WithCause(err)
	}

	return nil
}

// BufferSize returns the size of the encoder's buffer.
func (e *streamEncoder) BufferSize() int {
	return e.bufferSize
}

// SetBufferSize allows changing the encoder's buffer size.
// It returns an error if the new size is invalid.
func (e *streamEncoder) SetBufferSize(size int) error {
	if size <= 0 {
		return NewJSONError(ErrInvalidOptions, "buffer size must be positive")
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.writer = bufio.NewWriterSize(e.writer, size)
	e.bufferSize = size

	return nil
}

// Verify interface implementation at compile time
var _ JSONEncoder = (*streamEncoder)(nil)
