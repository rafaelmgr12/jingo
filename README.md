# JSON Parser in Go

A recursive descent JSON parser implemented in Go that supports parsing JSON objects and arrays with various data types. This parser follows the JSON specification and provides detailed error reporting with line and column information.

## Features

- Full JSON specification support including:
  - Objects and nested objects
  - Arrays and nested arrays
  - String values
  - Number values (including negative and decimal numbers)
  - Boolean values (true/false)
  - Null values
- Detailed error reporting with line and column numbers
- Lexical analysis with proper token recognition
- Abstract Syntax Tree (AST) generation
- Whitespace and newline handling
- Support for escape sequences in strings

## Project Structure

```bash
jsongoparser/
├── token.go       // Token types and definitions
├── lexer.go       // Lexical analyzer
├── interface.go   // AST interfaces
├── ast.go         // AST node implementations
├── parser.go      // JSON parser
└── parser_test.go // Parser tests
```

## Components

### Token

Defines the token types and structures used in JSON parsing:

- Structural tokens ({, }, [, ], :, ,)
- Value tokens (string, number, true, false, null)
- Special tokens (EOF, ILLEGAL)

### Lexer

Performs lexical analysis of the input JSON string:

- Converts input into a stream of tokens
- Handles whitespace and newlines
- Tracks line and column numbers
- Supports all JSON value types

### Parser

Performs syntactic analysis and builds an Abstract Syntax Tree:

- Recursive descent parsing
- Proper error handling
- Support for nested structures
- Validates JSON syntax

### AST (Abstract Syntax Tree)

Represents the structure of the JSON data:

- Objects with key-value pairs
- Arrays with ordered elements
- Various literal types (string, number, boolean, null)

## Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/rafaelmgr12/jsongoparser"
)

func main() {
    // Create a new lexer with JSON input
    input := `{"name": "John", "age": 30}`
    lexer := jsongoparser.NewLexer(input)

    // Create a parser with the lexer
    parser := jsongoparser.NewParser(lexer)

    // Parse the JSON and handle any errors
    value, err := parser.ParseJSON()
    if err != nil {
        log.Fatalf("Error parsing JSON: %v", err)
    }

    // Type assert to access the parsed data
    if obj, ok := value.(*jsongoparser.Object); ok {
        // Access object properties
        fmt.Println(obj.Pairs)
    }
}
```

## Error Handling

The parser provides detailed error messages including line and column information:

```go
package main

import (
    "fmt"
    "github.com/rafaelmgr12/jsongoparser"
)

func main() {
    input := `{"name": "John", age: 30}`  // Missing quotes around 'age'
    lexer := jsongoparser.NewLexer(input)
    parser := jsongoparser.NewParser(lexer)

    value, err := parser.ParseJSON()
    if err != nil {
        fmt.Println(err)
        // Output: Line 1, Column 14: expected string key, got age
    }
}
```

## Running Tests

To run the test suite:

```bash
go test ./...
```

## Known Issues and Limitations

- **Number Handling**: Currently, number values are stored as strings and need to be parsed when used. You may encounter issues when performing numerical operations directly on these values.
- **Error Recovery**: The parser could benefit from improved error recovery mechanisms to handle and report multiple errors in a user-friendly manner.
- **String Representations**: The `String()` methods are simplified and may not provide complete representations of complex JSON structures.
- **Performance**: Large JSON files may not be handled efficiently, and performance optimizations could be introduced.
- **Number Validation**: Additional validation could be added for number formats to ensure compliance with the JSON specification.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
