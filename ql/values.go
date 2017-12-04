package ql

type Value interface {
	valueType()
}

// EnumValue is the values of some enum type.
type EnumValue struct {
	name  string
	value int
}

func (e *EnumValue) valueType() {}
