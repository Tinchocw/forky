package interpreter

import "sync"

type Env struct {
	variables map[string]Value
	parent    *Env
	mu        *sync.RWMutex
}

func NewEnv(parent *Env) *Env {
	return &Env{
		variables: make(map[string]Value),
		parent:    parent,
		mu:        &sync.RWMutex{},
	}
}

func (e *Env) GetVariable(name string) (Value, bool) {
	e.mu.RLock()
	val, ok := e.variables[name]
	e.mu.RUnlock()
	if !ok && e.parent != nil {
		return e.parent.GetVariable(name)
	}
	return val, ok
}

func (e *Env) DefineVariable(name string, val Value) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.variables[name]; exists {
		return false
	}
	e.variables[name] = val
	return true
}

func (e *Env) AssignVariable(name string, val Value) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.variables[name]; exists {
		e.variables[name] = val
		return true
	}
	if e.parent != nil {
		return e.parent.AssignVariable(name, val)
	}
	return false
}
