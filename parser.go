// parser.go
package jsongoparser

import (
	"fmt"
)

// Parser holds the state while parsing JSON input. It maintains the current token and the next token,
// along with a list of any errors encountered during parsing.
type Parser struct {
	// lexer provides tokens from the input string.
	lexer *Lexer
	// currentToken is the current token being examined.
	currentToken Token
	// peekToken is the next token in the stream.
	peekToken Token
	// errors is a collection of parsing errors.
	errors []string
}

// NewParser creates a new Parser instance for the given lexer.
//
// The function initializes the Parser by reading two tokens
// to set up the currentToken and peekToken fields.
func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	// Read two tokens to initialize currentToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances to the next token in the token stream.
// It updates currentToken to the value of peekToken,
// and then gets a new value for peekToken from the lexer.
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// ParseJSON is the entry point for parsing JSON content. It returns the parsed
// Value and an error if the parsing fails.
//
// The function expects the JSON input to start with either a '{' or a '['.
func (p *Parser) ParseJSON() (Value, error) {
	var value Value

	// JSON must start with either { or [
	switch p.currentToken.Type {
	case TokenBraceOpen:
		value = p.parseObject()
	case TokenBracketOpen:
		value = p.parseArray()
	default:
		return nil, fmt.Errorf("expected { or [, got %s at line %d, column %d",
			p.currentToken.Type, p.currentToken.Line, p.currentToken.Column)
	}

	// After parsing the main value, we should be at EOF
	if p.peekToken.Type != TokenEOF {
		return nil, fmt.Errorf("unexpected token after main value: %s at line %d, column %d",
			p.peekToken.Type, p.peekToken.Line, p.peekToken.Column)
	}

	return value, nil
}

// parseObject parses a JSON object: { "key": value, ... }.
// It returns an Object value containing the key-value pairs.
func (p *Parser) parseObject() Value {
	object := &Object{
		Token: p.currentToken, // Store opening {
		Pairs: make(map[string]Value),
	}

	// Handle empty object case: {}
	if p.peekToken.Type == TokenBraceClose {
		p.nextToken()
		return object
	}

	p.nextToken() // move past {

	// Parse first key-value pair
	key, value := p.parseKeyValuePair()
	object.Pairs[key] = value

	// Parse additional key-value pairs
	for p.peekToken.Type == TokenComma {
		p.nextToken() // move past comma
		p.nextToken() // move to next key
		key, value = p.parseKeyValuePair()
		object.Pairs[key] = value
	}

	// Ensure we have a closing }
	if p.peekToken.Type != TokenBraceClose {
		p.addError("expected }, got %s", p.peekToken.Type)
		return nil
	}

	p.nextToken() // move past }
	return object
}

// parseKeyValuePair parses a key-value pair in a JSON object.
// It returns the key as a string and the value as a Value.
func (p *Parser) parseKeyValuePair() (string, Value) {
	// Key must be a string
	if p.currentToken.Type != TokenString {
		p.addError("expected string key, got %s", p.currentToken.Type)
		return "", nil
	}

	key := p.currentToken.Literal

	// Must have a colon after key
	if p.peekToken.Type != TokenColon {
		p.addError("expected :, got %s", p.peekToken.Type)
		return "", nil
	}

	p.nextToken() // move past key
	p.nextToken() // move past colon

	value := p.parseValue()

	return key, value
}

// parseArray parses a JSON array: [ value, value, ... ].
// It returns an Array value containing the elements.
func (p *Parser) parseArray() Value {
	array := &Array{
		Token:    p.currentToken,
		Elements: []Value{},
	}

	// Handle empty array case: []
	if p.peekToken.Type == TokenBracketClose {
		p.nextToken()
		return array
	}

	p.nextToken() // move past [

	// Parse first value
	value := p.parseValue()
	array.Elements = append(array.Elements, value)

	// Parse additional values
	for p.peekToken.Type == TokenComma {
		p.nextToken() // move past comma
		p.nextToken() // move to next value
		value = p.parseValue()
		array.Elements = append(array.Elements, value)
	}

	// Ensure we have a closing ]
	if p.peekToken.Type != TokenBracketClose {
		p.addError("expected ], got %s", p.peekToken.Type)
		return nil
	}

	p.nextToken() // move past ]
	return array
}

// parseValue parses any JSON value. It returns the parsed value.
//
// The function handles strings, numbers, booleans, nulls, objects, and arrays.
func (p *Parser) parseValue() Value {
	switch p.currentToken.Type {
	case TokenString:
		return &StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}

	case TokenNumber:
		return &NumberLiteral{Token: p.currentToken, Value: p.currentToken.Literal}

	case TokenTrue:
		return &Boolean{Token: p.currentToken, Value: true}

	case TokenFalse:
		return &Boolean{Token: p.currentToken, Value: false}

	case TokenNull:
		return &Null{Token: p.currentToken}

	case TokenBraceOpen:
		return p.parseObject()

	case TokenBracketOpen:
		return p.parseArray()

	default:
		p.addError("unexpected token %s", p.currentToken.Type)
		return nil
	}
}

// addError adds a formatted error message to the parser's error list.
//
// The function records the error message along with the line and column numbers
// where the error occurred.
func (p *Parser) addError(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	p.errors = append(p.errors, fmt.Sprintf("Line %d, Column %d: %s",
		p.currentToken.Line, p.currentToken.Column, msg))
}

// Errors returns all parsing errors encountered by the parser.
func (p *Parser) Errors() []string {
	return p.errors
}
