package examples

import (
	"fmt"
	"strconv"

	"github.com/rafaelmgr12/jingo/pkg/parser"
)

func Example() {
	input := `{
        "name": "John Doe",
        "age": 30,
        "address": {
            "street": "123 Main St",
            "city": "New York"
        }
    }`

	lexer := parser.NewLexer(input)
	p := parser.NewParser(lexer)

	value, err := p.ParseJSON()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	obj := value.(*parser.Object)
	fmt.Printf("Name: %s\n", obj.Pairs["name"].(*parser.StringLiteral).Value)
	age, _ := strconv.Atoi(obj.Pairs["age"].(*parser.NumberLiteral).Value)
	fmt.Printf("Age: %d\n", age)
}

func ExampleParser_ParseJSON() {
	input := `{"key": "value"}`
	lexer := parser.NewLexer(input)
	p := parser.NewParser(lexer)

	value, err := p.ParseJSON()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	obj := value.(*parser.Object)
	fmt.Printf("Value: %s\n", obj.Pairs["key"].(*parser.StringLiteral).Value)
	// Output: Value: value
}
