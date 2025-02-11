package encoding

// Marshaler is the interface implemented by types that can marshal themselves into valid JSON.
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

// Unmarshaler is the interface implemented by types that can unmarshal a JSON description of themselves.
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}
