package ql

// Executor for executing GraphQL requests.
type Executor struct{}

func (e *Executor) executeRequest(schema Schema, doc Document, oper string, vals map[string]string, initial Object) {
}

func (e *Executor) executeQuery()                     {}
func (e *Executor) executeMutation()                  {}
func (e *Executor) getOperation()                     {}
func (e *Executor) coerceVariableValues()             {}
func (e *Executor) executeSelectionSet(isNormal bool) {}
func (e *Executor) collectFields()                    {}
func (e *Executor) doesFragmentTypeApply()            {}
func (e *Executor) executeField()                     {}
func (e *Executor) coerceArgumentValues()             {}
func (e *Executor) resolveFieldValue()                {}
func (e *Executor) completeValue()                    {}
func (e *Executor) resolveAbstractValue()             {}
func (e *Executor) mergeSelectionSets()               {}
