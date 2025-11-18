package interpreter

type Env struct {
	variables map[string]*Value
	parent    *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{
		variables: make(map[string]*Value),
		parent:    parent,
	}
}

func (e *Env) GetVariable(name string) (*Value, bool) {
	val, ok := e.variables[name]
	if !ok && e.parent != nil {
		return e.parent.GetVariable(name)
	}
	return val, ok
}

func (e *Env) DefineVariable(name string, val Value) bool {
	if _, exists := e.variables[name]; exists {
		return false
	}
	e.variables[name] = &val
	return true
}

func (e *Env) AssignVariable(name string, val Value) bool {
	if _, exists := e.variables[name]; exists {
		e.variables[name] = &val
		return true
	}
	if e.parent != nil {
		return e.parent.AssignVariable(name, val)
	}
	return false
}
