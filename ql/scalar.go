package ql

// // Info for type name and description.
// type Info struct {
// 	name string
// 	desc string
// }
//
// // Meta information for type system.
// type Meta struct {
// 	Info
// 	kind string
// }
//
// // ScalarMeta represents meta information for scalar types.
// type ScalarMeta struct {
// 	Meta
// }
//
// // Scalar defines for all scalar types: Int, Float, String, Boolean, ID and
// // custom-defined.
// type Scalar struct {
// 	meta   ScalarMeta
// 	result CoerceFunc
// 	input  CoerceFunc
// }
//
// // CoerceFunc is the function type for serialize and deserialize.
// type CoerceFunc func(value interface{}) interface{}

// Document is a type-system type parsed from ast.Document.
type Document struct {
	opers []Operation
}

// Schema is a type-system type parsed from ast.SchemaDefinition.
type Schema struct {
	RootQuery    Object
	RootMutation Object
}

// Object .
type Object struct{}

// Interface .
type Interface struct{}

// Union .
type Union struct{}

// Operation .
type Operation struct {
	operType string
	varDefns []VariableDefinition
	selset   SelectionSet
}

// GroupedFieldSet .
type GroupedFieldSet struct {
	rspKey string
	fields []Field
}

// Field .
type Field struct {
	name string
	args []Argument
}

// Argument .
type Argument struct{}

// Selection is an interface for Field, FragmentSpread, InlineFragment
type Selection interface{}

// SelectionSet .
type SelectionSet struct{}

// VariableDefinition .
type VariableDefinition struct {
	name    string
	typ     string // TODO
	deflVal string
}

// Inputer for input coercsion.
type Inputer interface {
	input(value interface{}) interface{}
}

// Type is the interface for all variable type.
type Type interface{}

// func resultInt(value interface{}) interface{} {
// 	switch value := value.(type) {
// 	case int:
// 		if value > math.MaxInt32 || value < math.MinInt32 {
// 			return errors.New("integer out of range(-2^31~2^31)")
// 		}
// 		return value
// 		// case int8:
// 		// case int16:
// 		// case int32:
// 		// case int64:
// 		// case uint:
// 		// case uint8:
// 		// case uint16:
// 		// case uint32:
// 		// case uint64:
// 		// case float32:
// 		// case float64:
// 		// case string:
// 		// case bool:
// 		// case *int:
// 		// case *int8:
// 		// case *int16:
// 		// case *int32:
// 		// case *int64:
// 		// case *uint:
// 		// case *uint8:
// 		// case *uint16:
// 		// case *uint32:
// 		// case *uint64:
// 		// case *float32:
// 		// case *float64:
// 		// case *string:
// 		// case *bool:
// 		// default:
// 	}
// 	return nil
// }

// func inputInt(value interface{}) interface{} {}
