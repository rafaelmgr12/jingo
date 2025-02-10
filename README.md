# Jingo

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
jingo/
├── pkg/
│   ├── parser/           # Core parsing components
│   │   ├── ast.go        # Abstract Syntax Tree implementation
│   │   ├── lexer.go      # Lexical analyzer
│   │   ├── parser.go     # JSON parser
│   │   ├── token.go      # Token definitions
│   │   └── interface.go  # Parser interfaces
│   └── encoding/         # Encoding/decoding layer
│       ├── json.go       # Main Marshal/Unmarshal implementation
│       └── stream.go     # Streaming encoder/decoder
├── examples/             # Usage examples
├── docs/                 # Documentation
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

### Parsing JSON

First, you need to parse JSON strings into Go structures. Here’s an example that demonstrates how to parse a JSON string:

```go
package main

import (
    "fmt"
    "log"
    "github.com/rafaelmgr12/jingo/pkg/parser"
)

func main() {
    // JSON input
    input := `{"name": "John", "age": 30}`
    lexer := parser.NewLexer(input)
    p := parser.NewParser(lexer)
    value, err := p.ParseJSON()
    if err != nil {
        log.Fatalf("Error parsing JSON: %v", err)
    }

    // Access parsed data
    if obj, ok := value.(*parser.Object); ok {
        fmt.Println("Parsed JSON:", obj.Pairs)
    }
}
```

This example creates a lexer and a parser, then parses the JSON string, and finally prints the parsed data.

### Serializing JSON

You can also serialize Go data structures back into JSON strings:

```go
package main

import (
    "fmt"
    "log"
    "github.com/rafaelmgr12/jingo/pkg/encoding"
)

func main() {
    // Go data structure to be serialized
    data := map[string]interface{}{
        "name": "John",
        "age":  30,
        "address": map[string]string{
            "street": "123 Main St",
            "city":   "New York",
        },
    }

    // Serialize the data to JSON
    jsonStr, err := encoding.Marshal(data)
    if err != nil {
        log.Fatalf("Error serializing to JSON: %v", err)
    }
    fmt.Println("Serialized JSON:", string(jsonStr))
}
```

This example shows how to use the `Marshal` function to convert a Go data structure into a JSON string.

### Sending JSON over HTTP

You can send JSON data over HTTP using the following example:

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/rafaelmgr12/jingo/pkg/encoding"
    "github.com/rafaelmgr12/jingo/pkg/parser"
)

func main() {
    // JSON input
    input := `{"name": "John Doe", "age": 30}`
    lexer := parser.NewLexer(input)
    p := parser.NewParser(lexer)
    value, err := p.ParseJSON()
    if err != nil {
        log.Fatalf("Error parsing JSON: %v", err)
    }

    // Serialize JSON
    jsonStr, err := value.(*parser.Object).ToJSON()
    if err != nil {
        log.Fatalf("Error serializing JSON: %v", err)
    }
    
    fmt.Println("Serialized JSON:", jsonStr)

    // Send JSON via HTTP
    headers := map[string]string{
        "Authorization": "Bearer example-token",
    }
    resp, err := SendJSON("http://example.com/api", jsonStr, headers)
    if err != nil {
        log.Fatalf("Error sending JSON: %v", err)
    }
    defer resp.Body.Close()

    fmt.Println("Response status:", resp.Status)
}

// SendJSON is a helper function to send JSON data over HTTP
func SendJSON(url, jsonStr string, headers map[string]string) (*http.Response, error) {
    client := &http.Client{}
    req, err := http.NewRequest("POST", url, strings.NewReader(jsonStr))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    for key, value := range headers {
        req.Header.Set(key, value)
    }
    return client.Do(req)
}
```

This example demonstrates how to parse JSON, serialize it, and then send it via an HTTP POST request.

## Error Handling

The parser provides detailed error messages including line and column information:

```go
package main

import (
    "fmt"
    "github.com/rafaelmgr12/jingo/pkg/parser"
)

func main() {
    input := `{"name": "John", age: 30}`  // Missing quotes around 'age'
    lexer := parser.NewLexer(input)
    p := parser.NewParser(lexer)

    value, err := p.ParseJSON()
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

- **Number Handling**:
  - Currently, number values are stored as both integers and floats as part of the `NumberLiteral` struct. This dual representation is cumbersome for direct numerical operations and requires explicit type checking and conversion by the user.
  
- **Error Recovery**:
  - The parser stops at the first encountered error. Improved error recovery mechanisms could be introduced to handle and report multiple errors gracefully, allowing partial parsing of valid sections of the JSON.

- **Streaming JSON**:
  - While the `readChunk` method exists to support streaming mode, its implementation needs thorough verification and testing to ensure it effectively handles large JSON documents streamed in chunks without missing or corrupting data.

- **Performance**:
  - Parsing large JSON files into memory could lead to inefficiencies, especially because the lexer and parser currently rely on in-memory strings and buffers. Optimizations could be made to improve performance, especially for memory-intensive operations.

- **String Representations**:
  - The `String()` methods are simplified and may not provide complete or accurate representations of complex JSON structures, particularly when handling nested objects or arrays with escape sequences.

- **Lack of Customization**:
  - While extensive, the existing configurations and error handling rules are somewhat rigid. Allowing more customization in terms of linting rules or parse-time options could enhance the utility of the parser for various use cases.

- **UTF-8 Handling**:
  - The lexer currently supports UTF-8 decoding, but there may be edge cases with complex Unicode characters or mixing different encodings which require thorough testing and validation to ensure robustness.

By addressing these issues, the JSON parser can become more robust, efficient, and user-friendly.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
