// stream.go
package encoding

import (
	"bufio"
	"io"
	"reflect"
	"sync"

	"github.com/rafaelmgr12/jingo/pkg/parser"
)

// streamDecoder provides a concrete implementation of JSONDecoder interface
type streamDecoder struct {
	reader     *bufio.Reader
	lexer      *parser.Lexer
	parser     *parser.Parser
	options    *Options
	mutex      sync.Mutex
	buffer     []byte
	bufferSize int // Added to track buffer size
}

// NewDecoder creates a new JSONDecoder implementation
func NewDecoder(r io.Reader, opts ...Option) (JSONDecoder, error) {
	options, err := applyOptions(opts...)
	if err != nil {
		return nil, NewJSONError(ErrInvalidOptions, "invalid decoder options").WithCause(err)
	}

	bufferSize := 4096
	if options.BufferSize > 0 {
		bufferSize = options.BufferSize
	}

	reader := bufio.NewReader(r)
	lexer := parser.NewLexer(reader)
	parser := parser.NewParser(lexer)

	return &streamDecoder{
		reader:     reader,
		lexer:      lexer,
		parser:     parser,
		options:    options,
		buffer:     make([]byte, bufferSize),
		bufferSize: bufferSize,
	}, nil
}

// Decode implements JSONDecoder.Decode
func (d *streamDecoder) Decode(v interface{}) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	value, err := d.parser.ParseJSON()
	if err != nil {
		return NewJSONError(ErrInvalidJSON, "failed to parse JSON stream").WithCause(err)
	}

	return unmarshalValue(value, reflect.ValueOf(v).Elem())
}

// More implements JSONDecoder.More
func (d *streamDecoder) More() bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	b, err := d.reader.Peek(1)
	if err != nil {
		return false
	}

	// Skip whitespace
	for len(b) > 0 && isWhitespace(b[0]) {
		if _, err := d.reader.ReadByte(); err != nil {
			return false
		}

		b, err = d.reader.Peek(1)
		if err != nil {
			return false
		}
	}

	return len(b) > 0
}

// BufferSize implements JSONDecoder.BufferSize
func (d *streamDecoder) BufferSize() int {
	return d.bufferSize
}

// isWhitespace helper function
func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

var _ JSONDecoder = (*streamDecoder)(nil)
