package jsongoparser

import (
	"fmt"
	"strconv"
	"strings"
)

// Object represents a JSON object - a collection of key-value pairs.
type Object struct {
	// Token is the opening '{' token
	Token Token
	// Pairs are the key-value pairs in the object.
	Pairs map[string]Value
}

// TokenLiteral returns the literal value of the token that defines the object.
func (o *Object) TokenLiteral() string { return o.Token.Literal }

// String returns a simplified string representation of the object.
func (o *Object) String() string {
	var b strings.Builder

	b.WriteString("{")

	i := 0
	for k, v := range o.Pairs {
		if i > 0 {
			b.WriteString(", ")
		}

		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(v.String())

		i++
	}

	b.WriteString("}")

	return b.String()
}

// valueNode is a placeholder method to ensure type safety within the Value interface.
func (o *Object) valueNode() {}

// Array represents a JSON array - an ordered list of values.
type Array struct {
	// Token is the opening '[' token.
	Token Token
	// Elements are the values in the array.
	Elements []Value
}

// TokenLiteral returns the literal value of the token that defines the array.
func (a *Array) TokenLiteral() string { return a.Token.Literal }

// String returns a simplified string representation of the array.
func (a *Array) String() string { return "[]" } // Simplified for now

// valueNode is a placeholder method to ensure type safety within the Value interface.
func (a *Array) valueNode() {}

// StringLiteral represents a JSON string value.
type StringLiteral struct {
	// Token is the string token.
	Token Token
	// Value is the actual string value.
	Value string
}

// TokenLiteral returns the literal value of the token that defines the string.
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }

// String returns the actual string value.
func (s *StringLiteral) String() string { return s.Value }

// valueNode is a placeholder method to ensure type safety within the Value interface.
func (s *StringLiteral) valueNode() {}

// NumberLiteral represents a JSON number value.
type NumberLiteral struct {
	// Token is the number token.
	Token Token
	// Value is the number as a string (we'll parse it when needed).
	Value string
	// Float is the actual float value of the number.
	Float float64
	// Int is the actual integer value of the number.
	Int int64
	// IsInt is a flag to indicate if the number is an integer.
	IsInt bool
	// IsValid is a flag to indicate if the number is valid JSON number.
	IsValid bool
}

// NewNumberLiteral creates a new NumberLiteral with proper validation and parsing
func NewNumberLiteral(token Token) *NumberLiteral {
	n := &NumberLiteral{
		Token: token,
		Value: token.Literal,
	}

	isInt := true // Assume it's an integer initially

	for i := 0; i < len(token.Literal); i++ {
		switch token.Literal[i] {
		case '-', '+':
			if i != 0 {
				// Signs should only be at the beginning
				return setInvalidNumberLiteral(n)
			}
		case '.':
			isInt = false
		case 'e', 'E':
			isInt = false
			// Ensure there's an exponent part
			if i+1 >= len(token.Literal) {
				return setInvalidNumberLiteral(n)
			}

			if token.Literal[i+1] == '-' || token.Literal[i+1] == '+' {
				i++ // Skip the sign in exponent
			}
		default:
			if token.Literal[i] < '0' || token.Literal[i] > '9' {
				return setInvalidNumberLiteral(n)
			}
		}
	}

	if isInt {
		i, err := strconv.ParseInt(token.Literal, 10, 64)
		if err != nil {
			return setInvalidNumberLiteral(n)
		}

		n.Int = i
		n.Float = float64(i)
	} else {
		f, err := strconv.ParseFloat(token.Literal, 64)
		if err != nil {
			return setInvalidNumberLiteral(n)
		}

		n.Float = f
	}

	n.IsValid = true
	n.IsInt = isInt

	return n
}

func setInvalidNumberLiteral(n *NumberLiteral) *NumberLiteral {
	n.IsValid = false
	n.IsInt = false
	n.Int = 0
	n.Float = 0

	return n
}

// TokenLiteral returns the literal value of the token that defines the number.
func (n *NumberLiteral) TokenLiteral() string { return n.Token.Literal }

// String returns the number value as a string.
func (n *NumberLiteral) String() string {
	if n.IsInt {
		return fmt.Sprintf("%d", n.Int)
	}

	return fmt.Sprintf("%f", n.Float)
}

// valueNode is a placeholder method to ensure type safety within the Value interface.
func (n *NumberLiteral) valueNode() {}

// IsValidNumber returns whether the number is a valid JSON number
func (n *NumberLiteral) IsValidNumber() bool {
	return n.IsValid
}

// Boolean represents a JSON boolean value (true or false).
type Boolean struct {
	// Token is the boolean token.
	Token Token
	// Value is the actual boolean value.
	Value bool
}

// TokenLiteral returns the literal value of the token that defines the boolean.
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

// String returns the boolean value as a string.
func (b *Boolean) String() string { return b.Token.Literal }

// valueNode is a placeholder method to ensure type safety within the Value interface.
func (b *Boolean) valueNode() {}

// Null represents a JSON null value.
type Null struct {
	// Token is the null token.
	Token Token
}

// TokenLiteral returns the literal value of the token that defines the null value.
func (n *Null) TokenLiteral() string { return n.Token.Literal }

// String returns the string representation of the null value.
func (n *Null) String() string { return "null" }

// valueNode is a placeholder method to ensure type safety within the Value interface.
func (n *Null) valueNode() {}
