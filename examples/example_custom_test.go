package examples

import (
	"fmt"

	"github.com/rafaelmgr12/jingo/pkg/encoding"
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
	fmt.Printf("UnmarshalJSON called with data: %s\n", data) // Debug

	var temp struct {
		CustomName string `json:"custom_name"`
		CustomAge  int    `json:"custom_age"`
	}

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

	// This output is used to validate the test
	// Output:
	// Marshaling Success: {"custom_name":"Alice","custom_age":28}
	// UnmarshalJSON called with data: {"custom_name":"Alice","custom_age":28}
	// Unmarshaling Success: {Name: Alice, Age: 28}
}
