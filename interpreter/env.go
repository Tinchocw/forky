package interpreter

type Env struct {
	variables map[string]Value
	// functions map[string]Function
	parent *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{
		variables: make(map[string]Value),
		// functions: make(map[string]Function),
		parent: parent,
	}
}

func (e *Env) GetVariable(name string) (Value, bool) {
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
	e.variables[name] = val
	return true
}

func (e *Env) AssignVariable(name string, val Value) bool {
	if _, exists := e.variables[name]; exists {
		e.variables[name] = val
		return true
	}
	if e.parent != nil {
		return e.parent.AssignVariable(name, val)
	}
	return false
}

// func (e *Env) GetFunction(name string) (Function, bool) {
// 	fn, ok := e.functions[name]
// 	if !ok && e.parent != nil {
// 		return e.parent.GetFunction(name)
// 	}
// 	return fn, ok
// }

// // TODO: define how the function scope works
// func (e *Env) DefineFunction(name string, fn Function) bool {
// 	if _, exists := e.functions[name]; exists {
// 		return false
// 	}
// 	e.functions[name] = fn
// 	return true
// }
