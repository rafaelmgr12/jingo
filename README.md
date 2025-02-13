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
- Custom marshaling and unmarshaling support through interfaces
- Streaming JSON encoding/decoding

## Project Structure

```bash
jingo/
├── pkg/
│   ├── parser/                   # Core parsing components
│   │   ├── ast.go                # Abstract Syntax Tree implementation
│   │   ├── lexer.go              # Lexical analyzer
│   │   ├── parser.go             # JSON parser
│   │   ├── token.go              # Token definitions
│   │   └── interface.go          # Parser interfaces
│   └── encoding/                 # Encoding/decoding layer
│       ├── json.go               # Main Marshal/Unmarshal implementation
│       ├── marshaller.go         # Marshaler and Unmarshaler interfaces
│       ├── options.go            # Configuration options
│       ├── stream_encoder.go     # Streaming encoder implementation
│       ├── stream_decoder.go     # Streaming decoder implementation
│       ├── interfaces.go         # Encoder/Decoder interfaces
│       └── errors.go             # Error definitions
├── examples/                     # Usage examples
│   ├── example_test.go           # General examples
│   ├── example_custom_test.go    # Custom marshaling/unmarshaling examples
├── docs/                         # Documentation
```

## Components

### Token

Defines the token types and structures used in JSON parsing:

- Structural tokens: (`{`, `}`, `[`, `]`, `:`, `,`)
- Value tokens: (string, number, true, false, null)
- Special tokens: (EOF, ILLEGAL)

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

### Encoding

Provides functions to marshal Go data structures into JSON strings and unmarshal JSON strings into Go data structures:

- `Marshal` and `Unmarshal` functions with optional configuration
- Support for custom marshaling/unmarshaling through interfaces (`Marshaler` and `Unmarshaler`)
- Options for controlling encoding/decoding behaviors, such as size limits and strict mode
- Streaming JSON encoder/decoder with buffer configurations

## Usage

### Parsing JSON

You need to parse JSON strings into Go structures. Here’s an example that demonstrates how to parse a JSON string:

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

### Custom Marshaling/Unmarshaling

You can define your own custom marshaling and unmarshaling for your types by implementing the `Marshaler` and `Unmarshaler` interfaces:

```go
package examples

import (
    "fmt"
    "github.com/rafaelmgr12/jingo/pkg/encoding"
    "testing"
)

// CustomStruct demonstrates a complex struct with custom JSON marshaling/unmarshaling
type CustomStruct struct {
    Name string
    Age  int
}

// MarshalJSON is a custom marshaling function
func (cs *CustomStruct) MarshalJSON() ([]byte, error) {
    return []byte(fmt.Sprintf(`{"custom_name":"%s","custom_age":%d}`, cs.Name, cs.Age)), nil
}

// UnmarshalJSON is a custom unmarshaling function
func (cs *CustomStruct) UnmarshalJSON(data []byte) error {
    var temp struct {
        CustomName string `json:"custom_name"`
        CustomAge  int    `json:"custom_age"`
    }
    fmt.Println("UnmarshalJSON called with data:", string(data))
    if err := encoding.Unmarshal(data, &temp); err != nil {
        return err
    }
    cs.Name = temp.CustomName
    cs.Age = temp.CustomAge
    return nil
}

func ExampleCustomStruct() {
    cs := &CustomStruct{Name: "Alice", Age: 28}

    // Test Marshaling
    data, err := encoding.Marshal(cs)
    if err != nil {
        fmt.Printf("Error marshaling custom struct: %v\n", err)
        return
    }

    expectedJSON := `{"custom_name":"Alice","custom_age":28}`
    gotJSON := string(data)
    if gotJSON != expectedJSON {
        fmt.Printf("Marshaling failed: expected %s, got %s\n", expectedJSON, gotJSON)
        return
    }
    fmt.Println("Marshaling Success:", gotJSON)

    // Test Unmarshaling
    newCS := &CustomStruct{}
    if err := encoding.Unmarshal([]byte(expectedJSON), newCS); err != nil {
        fmt.Printf("Error unmarshaling custom struct: %v\n", err)
        return
    }
    if newCS.Name != "Alice" || newCS.Age != 28 {
        fmt.Printf("Unmarshaling failed: expected {Name: Alice, Age: 28}, got {Name: %s, Age: %d}\n", newCS.Name, newCS.Age)
        return
    }
    fmt.Printf("Unmarshaling Success: {Name: %s, Age: %d}\n", newCS.Name, newCS.Age)

    // Output:
    // Marshaling Success: {"custom_name":"Alice","custom_age":28}
    // UnmarshalJSON called with data: {"custom_name":"Alice","custom_age":28}
    // Unmarshaling Success: {Name: Alice, Age: 28}
}
```

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
    jsonStr, err := encoding.Marshal(value)
    if err != nil {
        log.Fatalf("Error serializing JSON: %v", err)
    }
    
    fmt.Println("Serialized JSON:", string(jsonStr))

    // Send JSON via HTTP
    headers := map[string]string{
        "Authorization": "Bearer example-token",
    }

    resp, err := SendJSON("http://example.com/api", string(jsonStr), headers)
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
  - While streaming mode is supported through `stream_encoder.go` and `stream_decoder.go`, its implementation requires thorough verification and testing to ensure it effectively handles large JSON documents streamed in chunks without missing or corrupting data.

- **Performance**:
  - Parsing large JSON files into memory may lead to inefficiencies, especially because the lexer and parser currently rely on in-memory strings and buffers. Optimizations could be made to improve performance, especially for memory-intensive operations.

- **String Representations**:
  - The `String()` methods are simplified and might not provide accurate representations of complex JSON structures, particularly when handling nested objects or arrays with escape sequences.

- **Lack of Customization**:
  - While extensive, the existing configurations and error handling rules could be considered somewhat rigid. Allowing more customization in terms of linting rules or parse-time options could enhance the parser's utility for various use cases.

- **UTF-8 Handling**:
  - The lexer currently supports UTF-8 decoding, but edge cases with complex Unicode characters or mixed encodings require thorough testing and validation to ensure robustness.

By addressing these issues, the JSON parser can become more robust, efficient, and user-friendly.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
