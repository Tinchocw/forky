package common

import (
	"fmt"
	"strings"
)

/*
	Program			-> 	Statements * EOF
	Statements		-> 	BlockStatement 			|
							IfStatement 		|
							WhileStatement 		|
							FunctionDef 		|
							VarDeclaration 		|
							Assignment 			|
							PrintStatement 		|
							Expression

	BlockStatement	-> '{' Statements * '}'
	IfStatement 	-> 'if' '(' Expression ')' BlockStatement
						( 'else' 'if' '(' Expression ')' BlockStatement )*
						( 'else' BlockStatement )?
	WhileStatement 	-> 'while' '(' Expression ')' BlockStatement
	FunctionDef 	-> 'func' IDENTIFIER '(' Parameters? ')' BlockStatement
	VarDeclaration 	-> 'var' IDENTIFIER ( '=' Expression )? ';'
	Assignment 		-> IDENTIFIER '=' Expression ';'
	PrintStatement 	-> 'print' '(' Expression ')' ';'
	Return 			-> 'return' Expression ';'
*/

// Helper function to generate indentation
func indent(level int) string {
	return strings.Repeat("  ", level)
}

// Helper function to create tree connectors
func stmtTreePrefix(level int, isLast bool) string {
	if level == 0 {
		return ""
	}

	prefix := strings.Repeat("│   ", level-1)
	if isLast {
		prefix += "└── "
	} else {
		prefix += "├── "
	}
	return prefix
}

// Helper function for child connectors
func stmtChildPrefix(level int, isLast bool) string {
	if level == 0 {
		return "│   "
	}

	prefix := strings.Repeat("│   ", level)
	if isLast {
		// Remove the last vertical bar for the last child's children
		prefix = strings.Repeat("│   ", level-1) + "    "
	}
	return prefix
}

type Program struct {
	Statements []Statement
}

func (p Program) Print(level int) {
	fmt.Printf("%s%sProgram%s\n", stmtTreePrefix(level, false), COLOR_CYAN, COLOR_RESET)
	for i, stmt := range p.Statements {
		isLast := i == len(p.Statements)-1
		fmt.Printf("%s", stmtChildPrefix(level, false))
		switch s := stmt.(type) {
		case Program:
			s.Print(level + 1)
		case BlockStatement:
			s.Print(level + 1)
		case IfStatement:
			s.Print(level + 1)
		case WhileStatement:
			s.Print(level + 1)
		case FunctionDef:
			s.Print(level + 1)
		case VarDeclaration:
			s.Print(level + 1)
		case Assignment:
			s.Print(level + 1)
		case PrintStatement:
			s.Print(level + 1)
		case Return:
			s.Print(level + 1)
		case BreakStatement:
			s.Print(level + 1)
		}
		_ = isLast // Use the variable to avoid unused warning for now
	}
}

type Statement interface{}

type BlockStatement struct {
	Statements []Statement
}

func (bs BlockStatement) Print(level int) {
	fmt.Printf("%s%sBlockStatement%s\n", indent(level), COLOR_MAGENTA, COLOR_RESET)
	for _, stmt := range bs.Statements {
		switch s := stmt.(type) {
		case Program:
			s.Print(level + 1)
		case BlockStatement:
			s.Print(level + 1)
		case IfStatement:
			s.Print(level + 1)
		case WhileStatement:
			s.Print(level + 1)
		case FunctionDef:
			s.Print(level + 1)
		case VarDeclaration:
			s.Print(level + 1)
		case Assignment:
			s.Print(level + 1)
		case PrintStatement:
			s.Print(level + 1)
		case Return:
			s.Print(level + 1)
		case BreakStatement:
			s.Print(level + 1)
		}
	}
}

type IfStatement struct {
	Condition Expression
	Body      BlockStatement
	Else      *IfStatement
}

func (ifs IfStatement) Print(level int) {
	fmt.Printf("%s%sIfStatement%s\n", indent(level), COLOR_BLUE, COLOR_RESET)

	fmt.Printf("%s%sCondition:%s\n", indent(level+1), COLOR_YELLOW, COLOR_RESET)
	ifs.Condition.Print(level + 2)

	fmt.Printf("%s%sBody:%s\n", indent(level+1), COLOR_YELLOW, COLOR_RESET)
	ifs.Body.Print(level + 2)

	if ifs.Else != nil {
		fmt.Printf("%s%sElse:%s\n", indent(level+1), COLOR_YELLOW, COLOR_RESET)
		ifs.Else.Print(level + 2)
	}
}

type WhileStatement struct {
	Condition Expression
	Body      BlockStatement
}

func (ws WhileStatement) Print(level int) {
	fmt.Printf("%s%sWhileStatement%s\n", indent(level), COLOR_BLUE, COLOR_RESET)

	fmt.Printf("%s%sCondition:%s\n", indent(level+1), COLOR_YELLOW, COLOR_RESET)
	ws.Condition.Print(level + 2)

	fmt.Printf("%s%sBody:%s\n", indent(level+1), COLOR_YELLOW, COLOR_RESET)
	ws.Body.Print(level + 2)
}

type FunctionDef struct {
	Name       string
	Parameters []string
	Body       BlockStatement
}

func (fd FunctionDef) Print(level int) {
	fmt.Printf("%s%sFunctionDef%s\n", indent(level), COLOR_GREEN, COLOR_RESET)

	fmt.Printf("%s%sName:%s %s%s%s\n", indent(level+1), COLOR_YELLOW, COLOR_RESET, COLOR_WHITE, fd.Name, COLOR_RESET)

	if len(fd.Parameters) > 0 {
		fmt.Printf("%s%sParameters:%s\n", indent(level+1), COLOR_YELLOW, COLOR_RESET)
		for _, param := range fd.Parameters {
			fmt.Printf("%s%s%s%s\n", indent(level+2), COLOR_WHITE, param, COLOR_RESET)
		}
	}

	fmt.Printf("%s%sBody:%s\n", indent(level+1), COLOR_YELLOW, COLOR_RESET)
	fd.Body.Print(level + 2)
}

type VarDeclaration struct {
	Name  string
	Value Expression
}

func (vd VarDeclaration) Print(level int) {
	fmt.Printf("%s%sVarDeclaration%s\n", indent(level), COLOR_GREEN, COLOR_RESET)

	fmt.Printf("%s%sName:%s %s%s%s\n", indent(level+1), COLOR_YELLOW, COLOR_RESET, COLOR_WHITE, vd.Name, COLOR_RESET)

	fmt.Printf("%s%sValue:%s\n", indent(level+1), COLOR_YELLOW, COLOR_RESET)
	vd.Value.Print(level + 2)
}

type Assignment struct {
	Name  string
	Value Expression
}

func (a Assignment) Print(level int) {
	fmt.Printf("%s%sAssignment%s\n", indent(level), COLOR_GREEN, COLOR_RESET)

	fmt.Printf("%s%sName:%s %s%s%s\n", indent(level+1), COLOR_YELLOW, COLOR_RESET, COLOR_WHITE, a.Name, COLOR_RESET)

	fmt.Printf("%s%sValue:%s\n", indent(level+1), COLOR_YELLOW, COLOR_RESET)
	a.Value.Print(level + 2)
}

type PrintStatement struct {
	Value Expression
}

func (ps PrintStatement) Print(level int) {
	fmt.Printf("%s%sPrintStatement%s\n", indent(level), COLOR_RED, COLOR_RESET)

	fmt.Printf("%s%sValue:%s\n", indent(level+1), COLOR_YELLOW, COLOR_RESET)
	ps.Value.Print(level + 2)
}

type Return struct {
	Value Expression
}

func (r Return) Print(level int) {
	fmt.Printf("%s%sReturn%s\n", indent(level), COLOR_RED, COLOR_RESET)

	fmt.Printf("%s%sValue:%s\n", indent(level+1), COLOR_YELLOW, COLOR_RESET)
	r.Value.Print(level + 2)
}

type BreakStatement struct{}

func (bs BreakStatement) Print(level int) {
	fmt.Printf("%s%sBreakStatement%s\n", indent(level), COLOR_RED, COLOR_RESET)
}
