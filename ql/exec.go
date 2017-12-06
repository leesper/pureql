package ql

import (
	"fmt"

	"github.com/leesper/pureql/ql/ast"
)

// Runtime represents runtime type info extract from schema.
type Runtime struct {
	Scalars   map[string]*Scalar
	Objects   map[string]*Object
	Ifaces    map[string]*Interface
	Unions    map[string]*Union
	Enums     map[string]*Enum
	InputObjs map[string]*InputObject
	Lists     map[string]*List
	NonNulls  map[string]*NonNull
}

// NewRuntime returns a new Runtime shipped with type infos from schema. It returns
// error if schema is invalid.
func NewRuntime(schema *Schema) (*Runtime, error) {
	if err := validateSchema(schema); err != nil {
		return nil, err
	}

	runtime := &Runtime{
		Scalars:   make(map[string]*Scalar),
		Objects:   make(map[string]*Object),
		Ifaces:    make(map[string]*Interface),
		Unions:    make(map[string]*Union),
		Enums:     make(map[string]*Enum),
		InputObjs: make(map[string]*InputObject),
		Lists:     make(map[string]*List),
		NonNulls:  make(map[string]*NonNull),
	}
	if schema == nil {
		return runtime, nil
	}
	extractObjectTypes(runtime, schema.Qry)
	extractObjectTypes(runtime, schema.Mut)
	return runtime, nil
}

func extractObjectTypes(runtime *Runtime, obj *Object) {
	if obj == nil {
		return
	}
	if _, ok := runtime.Objects[obj.Name]; ok {
		return
	}

	runtime.Objects[obj.Name] = obj
	for _, iface := range obj.Ifaces {
		extractIfaceTypes(runtime, iface)
	}
	for _, field := range obj.Fields {
		extractFieldTypes(runtime, field)
	}
}

func extractIfaceTypes(runtime *Runtime, iface *Interface) {
	if iface == nil {
		return
	}
	if _, ok := runtime.Ifaces[iface.Name]; ok {
		return
	}

	runtime.Ifaces[iface.Name] = iface
	for _, field := range iface.Fields {
		extractFieldTypes(runtime, field)
	}
}

func extractInputObjectTypes(runtime *Runtime, iobj *InputObject) {
	if iobj == nil {
		return
	}
	if _, ok := runtime.InputObjs[iobj.Name]; ok {
		return
	}

	runtime.InputObjs[iobj.Name] = iobj
	for _, field := range iobj.Fields {
		extractFieldTypes(runtime, field)
	}
}

func extractFieldTypes(runtime *Runtime, field *Field) {
	if field == nil {
		return
	}
	extractTypes(runtime, field.Typ)
}

func extractListTypes(runtime *Runtime, list *List) {
	if list == nil {
		return
	}
	if _, ok := runtime.Lists[list.Type()]; ok {
		return
	}

	runtime.Lists[list.Type()] = list
	extractTypes(runtime, list.OfType)
}

func extractNonNullTypes(runtime *Runtime, nn *NonNull) {
	if nn == nil {
		return
	}
	if _, ok := runtime.NonNulls[nn.Type()]; ok {
		return
	}

	runtime.NonNulls[nn.Type()] = nn
	extractTypes(runtime, nn.OfType)
}

func extractUnionTypes(runtime *Runtime, union *Union) {
	if union == nil {
		return
	}
	if _, ok := runtime.Unions[union.Name]; ok {
		return
	}
	for _, typ := range union.Typs {
		extractTypes(runtime, typ)
	}
}

func extractTypes(runtime *Runtime, typ Type) {
	switch typ := typ.(type) {
	case *Enum:
		runtime.Enums[typ.Name] = typ
	case *InputObject:
		extractInputObjectTypes(runtime, typ)
	case *Interface:
		extractIfaceTypes(runtime, typ)
	case *List:
		extractListTypes(runtime, typ)
	case *NonNull:
		extractNonNullTypes(runtime, typ)
	case *Object:
		extractObjectTypes(runtime, typ)
	case *Scalar:
		runtime.Scalars[typ.Name] = typ
	case *Union:
		extractUnionTypes(runtime, typ)
	default:
		panic(fmt.Errorf("unexpected type %T", typ))
	}
}

func (runtime *Runtime) findType(name string) Type {
	var typ Type
	var ok bool
	if typ, ok = runtime.Scalars[name]; ok {
		return typ
	}
	if typ, ok = runtime.Objects[name]; ok {
		return typ
	}
	if typ, ok = runtime.Ifaces[name]; ok {
		return typ
	}
	if typ, ok = runtime.Unions[name]; ok {
		return typ
	}
	if typ, ok = runtime.Enums[name]; ok {
		return typ
	}
	if typ, ok = runtime.InputObjs[name]; ok {
		return typ
	}
	if typ, ok = runtime.Lists[name]; ok {
		return typ
	}
	if typ, ok = runtime.NonNulls[name]; ok {
		return typ
	}
	return nil
}

// Response of executing request.
type Response struct {
	Errors []error
}

// Execute executes the request defined by document with optional variable values.
func (runtime *Runtime) Execute(document *ast.Document, operationName string, variableValues map[string]interface{}) *Response {
	rsp := &Response{}

	// TODO
	// err = validateDocument(document)
	// if err != nil {
	// 	rsp.Errors = append(rsp.Errors, err)
	// 	return rsp
	// }

	operation, err := runtime.getOperation(document, operationName)
	if err != nil {
		rsp.Errors = append(rsp.Errors, err)
		return rsp
	}

	coercedVarVals, err := runtime.coerceVariableValues(operation, variableValues)
	if err != nil {
		rsp.Errors = append(rsp.Errors, err)
		return rsp
	}
	return runtime.executeRequest(operation, coercedVarVals)
}

func (runtime *Runtime) getOperation(document *ast.Document, operationName string) (*ast.OperationDefinition, error) {
	var oper *ast.OperationDefinition
	var ok bool
	if operationName == "" {
		count := 0
		for _, def := range document.Defs {
			oper, ok = def.(*ast.OperationDefinition)
			if ok {
				count++
			}
		}
		if count == 1 {
			return oper, nil
		}
		return nil, fmt.Errorf("query error: requiring operation name")

	}
	for _, def := range document.Defs {
		oper, ok = def.(*ast.OperationDefinition)
		if ok {
			if oper.Name.Text == operationName {
				return oper, nil
			}
		}
	}
	return nil, fmt.Errorf("query error: operation %s not found", operationName)
}

func (runtime *Runtime) coerceVariableValues(operation *ast.OperationDefinition, variableValues map[string]interface{}) (map[string]interface{}, error) {
	coercedValues := map[string]interface{}{}
	for _, varDefn := range operation.VarDefns.VarDefns {
		varName := varDefn.Var.Name.Text
		varType := runtime.resolveASTType(varDefn.Typ)
		deflVal := varDefn.DeflVal.Val
		value, ok := variableValues[varName]
		if !ok {
			if deflVal != nil {
				coercedValues[varName] = deflVal
			} else if isNonNull(varType) {
				return nil, fmt.Errorf("query error: type %T is non-null", varType)
			}
		} else {
			// coerce value
			coercedVal, err := coerceInputValue(varType, value)
			if err != nil {
				return nil, err
			}
			coercedValues[varName] = coercedVal
		}
	}
	return coercedValues, nil
}

func coerceInputValue(typ Type, value interface{}) (interface{}, error) {
	switch typ := typ.(type) {
	case *Object:
		return nil, fmt.Errorf("invalid input object %s", typ.Name)
	case *Interface:
		return nil, fmt.Errorf("invalid input interface %s", typ.Name)
	case *Union:
		return nil, fmt.Errorf("invalid input union %s", typ.Name)
	case *Scalar:
		return nil, fmt.Errorf("invalid input type %T", typ)
	case *Enum:
		return nil, fmt.Errorf("invalid input type %T", typ)
	case *InputObject:
		return nil, fmt.Errorf("invalid input type %T", typ)
	case *List:
		return nil, fmt.Errorf("invalid input type %T", typ)
	case *NonNull:
		return nil, fmt.Errorf("invalid input type %T", typ)
	default:
		return nil, fmt.Errorf("invalid input type %T", typ)
	}
}

func (runtime *Runtime) executeRequest(operation *ast.OperationDefinition, coercedVariableValues map[string]interface{}) *Response {
	return nil
}

func isNonNull(typ Type) bool {
	_, ok := typ.(*NonNull)
	return ok
}

func (runtime *Runtime) resolveASTType(astTyp ast.Type) Type {
	return runtime.findType(formatName(astTyp))
}

func formatName(typ ast.Type) string {
	var bang string
	switch typ := typ.(type) {
	case *ast.NamedType:
		if typ.NonNull {
			bang = "!"
		}
		return typ.Name.Text + bang
	case *ast.ListType:
		if typ.NonNull {
			bang = "!"
		}
		return fmt.Sprintf("[ %s ]", formatName(typ.Typ)) + bang
	default:
		panic(fmt.Errorf("unexpected AST type %T", typ))
	}
}
