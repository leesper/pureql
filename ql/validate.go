package ql

import (
	"fmt"
	"reflect"
)

func ruleMustDefineOneOrMoreFields(typ Type) error {
	var numOfFields = 0
	switch typ := typ.(type) {
	case *Object:
		numOfFields = len(typ.Fields)
	case *Interface:
		numOfFields = len(typ.Fields)
	case *InputObject:
		numOfFields = len(typ.Fields)
	default:
		return fmt.Errorf("type %T not applied to this rule", typ)
	}

	if numOfFields <= 0 {
		return fmt.Errorf("type %T must define one or more fields", typ)
	}
	return nil
}

func ruleFieldsMustHaveUniqueNamesWithin(typ Type) error {
	var fields []*Field
	switch typ := typ.(type) {
	case *Object:
		fields = typ.Fields
	case *Interface:
		fields = typ.Fields
	case *InputObject:
		fields = typ.Fields
	default:
		return fmt.Errorf("type %T not applied to this rule", typ)
	}

	fieldCount := map[string]int{}
	for _, f := range fields {
		fieldCount[f.Name]++
		if fieldCount[f.Name] > 1 {
			return fmt.Errorf("type %T has multiple fields named %s", typ, f.Name)
		}
	}
	return nil
}

func ruleMustBeSuperSetOfAllIfaces(obj *Object) error {
	for _, iface := range obj.Ifaces {
		if err := ruleMustIncludeFieldOfSameName(obj, iface); err != nil {
			return err
		}
	}
	return nil
}

func ruleMustIncludeFieldOfSameName(obj *Object, iface *Interface) error {
	fieldMap := map[string]*Field{}
	for _, f := range obj.Fields {
		fieldMap[f.Name] = f
	}

	for _, f := range iface.Fields {
		if fieldMap[f.Name] == nil {
			return fmt.Errorf("object %s has no field %s of interface %s", obj.Name, f.Name, iface.Name)
		}
		err := ruleMustBeEqualOrSubTypeOf(fieldMap[f.Name].Typ, f.Typ)
		if err != nil {
			return err
		}
		err = ruleMustIncludeAgrumentOfSameName(fieldMap[f.Name].Defs, f.Defs)
		if err != nil {
			return err
		}
	}
	return nil
}

func ruleMustBeEqualOrSubTypeOf(typ, super Type) error {
	if reflect.DeepEqual(typ, super) {
		return nil
	}

	obj, isObject := typ.(*Object)
	iface, isIface := super.(*Interface)
	if isObject && isIface {
		return ruleMustIncludeFieldOfSameName(obj, iface)
	}

	union, isUnion := super.(*Union)
	if isObject && isUnion {
		for _, typ := range union.Typs {
			if err := ruleMustBeEqualOrSubTypeOf(obj, typ); err == nil {
				return nil
			}
		}
		return fmt.Errorf("type %T is not a sub-type of %T", typ, super)
	}

	lto, isListObject := typ.(*List)
	lti, isListIface := super.(*List)
	if isListObject && isListIface {
		return ruleMustBeEqualOrSubTypeOf(lto.OfType, lti.OfType)
	}

	nt, isNonNull := typ.(*NonNull)
	if isNonNull {
		return ruleMustBeEqualOrSubTypeOf(nt.OfType, super)
	}

	return nil
}

func ruleMustIncludeAgrumentOfSameName(args []*ArgDef, iargs []*ArgDef) error {
	argMap := map[string]*ArgDef{}
	for _, a := range args {
		argMap[a.name] = a
	}

	for _, a := range iargs {
		if argMap[a.name] == nil {
			return fmt.Errorf("no argument definition found for %s", a.name)
		}
		if !reflect.DeepEqual(argMap[a.name].typ, a.typ) {
			return fmt.Errorf("expected argument type %T, found %T", a.typ, argMap[a.name].typ)
		}
	}
	return nil
}

func ruleMustBeAllObjectTypes(uni *Union) error {
	for _, typ := range uni.Typs {
		if _, ok := typ.(*Object); !ok {
			return fmt.Errorf("type %T not an object type", typ)
		}
	}
	return nil
}

func ruleMustDefineOneOrMoreMemberTypes(uni *Union) error {
	if len(uni.Typs) <= 0 {
		return fmt.Errorf("union must define at least one type")
	}
	return nil
}

func ruleFieldOfInputObjectMustBeInputType(io *InputObject) error {
	for _, f := range io.Fields {
		switch f.Typ.(type) {
		case *Scalar, *Enum, *InputObject:
			// do nothing
		default:
			return fmt.Errorf("unexpected type %T, input type wanted", f.Typ)
		}
	}
	return nil
}
