package jsongoparser

// TokenType is a custom type that represents the type of a token in a JSON document. It is defined as
// a string type with a set of predefined values. Each value represents a different type of token that
// can be found in a JSON document.
type TokenType string

const (
	TokenBraceOpen    TokenType = "{"
	TokenBraceClose   TokenType = "}"
	TokenBracketOpen  TokenType = "["
	TokenBracketClose TokenType = "]"
	TokenColon        TokenType = ":"
	TokenComma        TokenType = ","
	TokenString       TokenType = "STRING"
	TokenNumber       TokenType = "NUMBER"
	TokenTrue         TokenType = "TRUE"
	TokenFalse        TokenType = "FALSE"
	TokenNull         TokenType = "NULL"
	TokenEOF          TokenType = "EOF"
	TokenIllegal      TokenType = "ILLEGAL"
)

// Token represents a token in a JSON document. It consists of a type, a literal value, and the line and
// column where the token was found in the document.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}
