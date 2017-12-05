package ql

// Field represents fields in Object, Interface and InputObject.
type Field struct {
	Name string
	Typ  interface{} // TODO: maybe a Type interface is appropriate
	Defs []*ArgDef
}

// ArgDef represents argument definitions in Object and Interface.
type ArgDef struct {
	name string
	typ  interface{}
	defl interface{}
	// directs []*Directive
}
