package ql

// Meta information for type system.
type Meta struct {
	kind string
	name string
	desc string
}

// ScalarMeta represents meta information for scalar types.
type ScalarMeta struct {
	Meta
}

// Scalar defines for all scalar types: Int, Float, String, Boolean, ID and
// custom-defined.
type Scalar struct {
	meta        ScalarMeta
	serialize   CoerceFunc
	deserialize CoerceFunc
}

// CoerceFunc is the function type for serialize and deserialize.
type CoerceFunc func(value interface{}) interface{}

// func serializeInt(value interface{}) interface{} {
//
// }
