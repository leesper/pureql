package ql

import "fmt"

// Executor for executing GraphQL requests.
type Executor struct{}

func (e *Executor) executeRequest(schema *Schema, doc *Document, operName string, vals map[string]string, initial *Object) interface{} {
	oper := e.getOperation(doc, operName)
	coercedValues := e.coerceVariableValues(schema, oper, vals)
	switch oper.operType {
	case query:
		return e.executeQuery(oper, schema, coercedValues, initial)
	case mutation:
		return e.executeMutation(oper, schema, coercedValues, initial)
	default:
		return fmt.Sprintf("unexpected type %s", oper.operType)
	}
}

func (e *Executor) getOperation(doc *Document, operName string) *Operation {
	return nil
}

func (e *Executor) coerceVariableValues(schema *Schema, oper *Operation, vals map[string]string) interface{} {
	return nil
}
func (e *Executor) executeQuery(oper *Operation, schema *Schema, vals interface{}, initial *Object) interface{} {
	return nil
}
func (e *Executor) executeMutation(oper *Operation, schema *Schema, vals interface{}, initial *Object) interface{} {
	return nil
}
func (e *Executor) executeSelectionSet(isNormal bool) {}
func (e *Executor) collectFields()                    {}
func (e *Executor) doesFragmentTypeApply()            {}
func (e *Executor) executeField()                     {}
func (e *Executor) coerceArgumentValues()             {}
func (e *Executor) resolveFieldValue()                {}
func (e *Executor) completeValue()                    {}
func (e *Executor) resolveAbstractValue()             {}
func (e *Executor) mergeSelectionSets()               {}
