package ql

// Schema is the entry point of GraphQL service.
type Schema struct {
	Qry *Object
	Mut *Object
}

// Scalar represents primitive value.
type Scalar struct{}

// Enum represents limited enumerable values.
type Enum struct{}

// Object defines a set of fields of another type in the type system.
type Object struct {
	Name   string
	Ifaces []*Interface
	Fields []*Field
}

// Interface defines an abstract type for Object to implement.
type Interface struct {
	Name   string
	Fields []*Field
}

// Union defines a list of possible Object types.
type Union struct {
	Name string
	Typs []interface{}
}

// List is a list of other types.
type List struct {
	OfType interface{}
}

// NonNull is a non-null wrapper of other types.
type NonNull struct {
	OfType interface{}
}

// InputObject is a struct for complex input.
type InputObject struct {
	Fields []*Field
}

// Directive represents directives the execution engine supports.
type Directive struct{}
