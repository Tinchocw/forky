package interpreter

import (
	"fmt"
	"sync"
)

type Env struct {
	variables *sync.Map
	parent    *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{
		variables: &sync.Map{},
		parent:    parent,
	}
}

func (e *Env) GetVariable(name string) (Value, error) {
	if val, ok := e.variables.Load(name); ok {
		return val.(Value), nil
	}
	if e.parent != nil {
		return e.parent.GetVariable(name)
	}
	return nil, fmt.Errorf("variable '%s' not defined", name)
}

func (e *Env) DefineVariable(name string, val Value) error {
	_, loaded := e.variables.LoadOrStore(name, val)
	if loaded {
		return fmt.Errorf("variable '%s' already defined in this scope", name)
	}
	return nil
}

func (e *Env) AssignVariable(name string, val Value) error {
	if _, ok := e.variables.Load(name); ok {
		e.variables.Store(name, val)
		return nil
	}
	if e.parent != nil {
		return e.parent.AssignVariable(name, val)
	}
	return fmt.Errorf("variable '%s' not defined", name)
}

func (e *Env) AssignArrayVariable(name string, indexes []int, val Value) error {
	if len(indexes) == 0 {
		return fmt.Errorf("no indexes provided")
	}

	arrayVal, err := e.GetVariable(name)
	if err != nil {
		return err
	}

	for i, index := range indexes {
		if arrayVal.Type() != VAL_ARRAY {
			return fmt.Errorf("variable '%s' is not an array", name)
		}
		av := arrayVal.(*ArrayValue)
		if index < 0 || index >= len(av.Values) {
			return fmt.Errorf("array index %d out of bounds", index)
		}

		if i == len(indexes)-1 {
			av.Values[index] = val
		} else {
			arrayVal = av.Values[index]
		}
	}

	return nil
}
