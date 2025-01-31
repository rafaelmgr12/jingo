package jsongoparser

// Node represents a node in our AST (Abstract Syntax Tree)
// Every node in our JSON structure will implement this interface
type Node interface {
	TokenLiteral() string // Returns the literal value of the token
	String() string       // Returns a string representation of the node
}

// Value represents any valid JSON value (string, number, object, array, bool, or null)
// This is also an interface since JSON can have different types of values
type Value interface {
	Node
	valueNode() // Dummy method to ensure type safety
}
