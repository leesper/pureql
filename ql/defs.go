package ql

// Field represents fields in Object, Interface and InputObject.
type Field struct {
	Name string
	Typ  Type
	Defs []*ArgDef
}

// ArgDef represents argument definitions in Object and Interface.
type ArgDef struct {
	Name string
	Typ  Type
	defl interface{}
	// directs []*Directive
}
