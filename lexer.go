package jsongoparser

import (
	"bufio"
	"io"
)

// Lexer is responsible for converting JSON input into a sequence of tokens.
// It maintains the current input string and tracks the positions of characters being read.
type Lexer struct {
	// The input string being tokenized.
	input string
	// The current position in the input (points to the current character).
	position int
	// The position in the input after the current character.
	readPosition int
	// The current character being examined.
	ch byte
	// The current line number in the input (1-based index).
	line int
	// The current column number in the input (0-based index).
	column int
	// The buffered reader for the input string.
	reader *bufio.Reader
	// The buffer is used to store read characters when streaming.
	buffer []byte
	// Flag to indicate if the lexer is in streaming mode.
	isStreaming bool
}

// NewLexer creates a new Lexer instance for the given input string.
func NewLexer(input interface{}) *Lexer {
	l := &Lexer{
		line:   1,
		column: 0,
		buffer: make([]byte, 4096),
	}

	switch v := input.(type) {
	case string:
		l.input = v
		l.isStreaming = false
	case io.Reader:
		l.reader = bufio.NewReader(v)
		l.isStreaming = true
		l.readChunk()
	}

	l.readChar()

	return l
}

// readChunk reads the next chunk of data from the input reader.
func (l *Lexer) readChunk() {
	if !l.isStreaming || l.reader == nil {
		return
	}

	remaining := len(l.input) - l.position
	if remaining > 0 {
		copy(l.buffer, l.input[l.position:])
		l.position += copy(l.buffer[remaining:], l.input[l.position:])

		return
	}

	l.input = ""

	n, err := l.reader.Read(l.buffer[remaining:])
	if n > 0 {
		l.input += string(l.buffer[remaining : remaining+n])
	}

	if err != nil && err != io.EOF {
		return
	}

	l.position = 0
	l.readPosition = 0
}

// NextToken retrieves the next token from the input, skipping any whitespace.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	currentLine := l.line
	currentColumn := l.column

	var t Token

	switch l.ch {
	case '{':
		t = Token{Type: TokenBraceOpen, Literal: string(l.ch), Line: currentLine, Column: currentColumn}
	case '}':
		t = Token{Type: TokenBraceClose, Literal: string(l.ch), Line: currentLine, Column: currentColumn}
	case '[':
		t = Token{Type: TokenBracketOpen, Literal: string(l.ch), Line: currentLine, Column: currentColumn}
	case ']':
		t = Token{Type: TokenBracketClose, Literal: string(l.ch), Line: currentLine, Column: currentColumn}
	case ':':
		t = Token{Type: TokenColon, Literal: string(l.ch), Line: currentLine, Column: currentColumn}
	case ',':
		t = Token{Type: TokenComma, Literal: string(l.ch), Line: currentLine, Column: currentColumn}
	case '"':
		return l.readString(currentLine, currentColumn)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		return l.readNumber(currentLine, currentColumn)
	case 't':
		return l.readTrue(currentLine, currentColumn)
	case 'f':
		return l.readFalse(currentLine, currentColumn)
	case 'n':
		return l.readNull(currentLine, currentColumn)
	case 0:
		t = Token{Type: TokenEOF, Literal: "", Line: currentLine, Column: currentColumn}
	default:
		t = Token{Type: TokenIllegal, Literal: string(l.ch), Line: currentLine, Column: currentColumn}
	}

	l.readChar()

	return t
}

// readChar advances the position in the input string and updates the current character.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		if l.isStreaming {
			l.readChunk()
		}

		if l.readPosition >= len(l.input) {
			l.ch = 0 // EOF
			return
		}
	}

	l.ch = l.input[l.readPosition]
	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// skipWhitespace skips over any whitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readString reads a string token.
func (l *Lexer) readString(line, column int) Token {
	var result []byte

	l.readChar()

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()

			if l.ch == 0 {
				return Token{Type: TokenIllegal, Literal: "Unterminated string", Line: line, Column: column}
			}

			result = append(result, '\\')
			result = append(result, l.ch)
		} else {
			result = append(result, l.ch)
		}

		l.readChar()
	}

	if l.ch == 0 {
		return Token{Type: TokenIllegal, Literal: "Unterminated string", Line: line, Column: column}
	}

	l.readChar()

	return Token{Type: TokenString, Literal: string(result), Line: line, Column: column}
}

// readNumber reads and validates a JSON number token.
func (l *Lexer) readNumber(line, column int) Token {
	start := l.position

	// Handle negative numbers
	if l.ch == '-' {
		l.readChar()

		if !isDigit(l.ch) {
			return Token{
				Type:    TokenIllegal,
				Literal: "Invalid number format: digit expected after '-'",
				Line:    line,
				Column:  column,
			}
		}
	}

	// First digit cannot be zero unless it's followed by a decimal point
	switch {
	case l.ch == '0':
		l.readChar()

		if isDigit(l.ch) {
			return Token{
				Type:    TokenIllegal,
				Literal: "Invalid number format: leading zeros not allowed",
				Line:    line,
				Column:  column,
			}
		}
	case isNonZeroDigit(l.ch):
		// Read integer part
		l.readChar()

		for isDigit(l.ch) {
			l.readChar()
		}
	case l.ch != '.': // If not a digit and not a decimal point, it's invalid
		return Token{
			Type:    TokenIllegal,
			Literal: "Invalid number format: expected digit",
			Line:    line,
			Column:  column,
		}
	}

	// Handle fractional part
	if l.ch == '.' {
		l.readChar()

		if !isDigit(l.ch) {
			return Token{
				Type:    TokenIllegal,
				Literal: "Invalid number format: digit expected after decimal point",
				Line:    line,
				Column:  column,
			}
		}

		for isDigit(l.ch) {
			l.readChar()
		}
	}

	// Handle exponential notation
	if l.ch == 'e' || l.ch == 'E' {
		l.readChar()

		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}

		if !isDigit(l.ch) {
			return Token{
				Type:    TokenIllegal,
				Literal: "Invalid number format: digit expected in exponent",
				Line:    line,
				Column:  column,
			}
		}

		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return Token{
		Type:    TokenNumber,
		Literal: l.input[start:l.position],
		Line:    line,
		Column:  column,
	}
}

// readTrue reads a true boolean token.
func (l *Lexer) readTrue(line, column int) Token {
	word := l.readWord()
	if word == "true" {
		return Token{Type: TokenTrue, Literal: "true", Line: line, Column: column}
	}

	return Token{Type: TokenIllegal, Literal: "Invalid token", Line: line, Column: column}
}

// readFalse reads a false boolean token.
func (l *Lexer) readFalse(line, column int) Token {
	word := l.readWord()
	if word == "false" {
		return Token{Type: TokenFalse, Literal: "false", Line: line, Column: column}
	}

	return Token{Type: TokenIllegal, Literal: "Invalid token", Line: line, Column: column}
}

// readNull reads a null token.
func (l *Lexer) readNull(line, column int) Token {
	word := l.readWord()
	if word == "null" {
		return Token{Type: TokenNull, Literal: "null", Line: line, Column: column}
	}

	return Token{Type: TokenIllegal, Literal: "Invalid token", Line: line, Column: column}
}

// readWord reads a word token (used for true, false, null).
func (l *Lexer) readWord() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// isLetter checks if a character is a letter.
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

// isDigit checks if a character is a digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isNonZeroDigit checks if a character is a non-zero digit.
func isNonZeroDigit(ch byte) bool {
	return '1' <= ch && ch <= '9'
}
