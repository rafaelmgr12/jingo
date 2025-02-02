package jsongoparser_test

import (
	"fmt"
	"strconv"

	"github.com/rafaelmgr12/jsongoparser"
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

	lexer := jsongoparser.NewLexer(input)
	parser := jsongoparser.NewParser(lexer)

	value, err := parser.ParseJSON()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	obj := value.(*jsongoparser.Object)
	fmt.Printf("Name: %s\n", obj.Pairs["name"].(*jsongoparser.StringLiteral).Value)
	age, _ := strconv.Atoi(obj.Pairs["age"].(*jsongoparser.NumberLiteral).Value)
	fmt.Printf("Age: %d\n", age)
}

func ExampleParser_ParseJSON() {
	input := `{"key": "value"}`
	lexer := jsongoparser.NewLexer(input)
	parser := jsongoparser.NewParser(lexer)

	value, err := parser.ParseJSON()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	obj := value.(*jsongoparser.Object)
	fmt.Printf("Value: %s\n", obj.Pairs["key"].(*jsongoparser.StringLiteral).Value)
	// Output: Value: value
}
