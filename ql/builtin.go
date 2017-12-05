package ql

// built-in scalar types.
var (
	Int     = &Scalar{Name: "Int"}
	Float   = &Scalar{Name: "Float"}
	String  = &Scalar{Name: "String"}
	Boolean = &Scalar{Name: "Boolean"}
	ID      = &Scalar{Name: "ID"}
)
