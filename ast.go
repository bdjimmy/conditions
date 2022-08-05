package condition

import (
	"bytes"
	"fmt"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ       ObjectType = "INTEGER"
	STRING_OBJ        ObjectType = "STRING"
	BOOLEAN_OBJ       ObjectType = "BOLLEAN"
	ARRAY_INTEGER_OBJ ObjectType = "ARRAY_INTEGER_OBJ"
	ARRAY_STRING_OBJ  ObjectType = "ARRAY_STRING_OBJ"
	FUNCTION_OBJ      ObjectType = "FUNCTION_OBJ"
	BUILTIN_OBJ       ObjectType = "BUILTIN_OBJ"
	NULL_OBJ          ObjectType = "NULL"
	ERROR_OBJ         ObjectType = "ERROR"
	// special
	IDENT_OBJ ObjectType = "IDENT_OBJ"
)

type Object interface {
	ObjectType() ObjectType
}

// Node 必须被每个AST阶段实现
type Node interface {
	node()
	String() string
}

// Expression 标识一个单一的表达式节点
type Expression interface {
	Node
	expressionNode()
}

// Program AST的root节点
type Program struct {
	Expression Expression
}

func (p *Program) node() {}
func (p *Program) String() string {
	var out bytes.Buffer
	out.WriteString(p.Expression.String())
	return out.String()
}

// Identifier 标识符字面量, abc bcd efg
type Identifier struct {
	Value string
}

func (i *Identifier) node()                  {}
func (i *Identifier) expressionNode()        {}
func (i *Identifier) ObjectType() ObjectType { return IDENT_OBJ }
func (i *Identifier) String() string         { return i.Value }

// Error 标识一个错误，用来传递
type Error struct {
	Message string
}

func (e *Error) ObjectType() ObjectType { return ERROR_OBJ }

// Integer 整形字面量, 123 456
type Integer struct {
	Value int64
}

func (il *Integer) node()                  {}
func (il *Integer) expressionNode()        {}
func (il *Integer) ObjectType() ObjectType { return INTEGER_OBJ }
func (il *Integer) String() string         { return fmt.Sprintf("%v", il.Value) }

// Boolean bool字面量, true false
type Boolean struct {
	Value bool
}

func (il *Boolean) node()                  {}
func (il *Boolean) expressionNode()        {}
func (il *Boolean) ObjectType() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) String() string          { return fmt.Sprintf("%v", b.Value) }

// String string字面量, "abc" "123"
type String struct {
	Value string
}

func (s *String) node()                  {}
func (s *String) expressionNode()        {}
func (s *String) ObjectType() ObjectType { return STRING_OBJ }
func (s *String) String() string         { return fmt.Sprintf("\"%s\"", s.Value) }

type ArrayInteger struct {
	Value []int64
}

func (a *ArrayInteger) node()                  {}
func (a *ArrayInteger) expressionNode()        {}
func (a *ArrayInteger) ObjectType() ObjectType { return ARRAY_INTEGER_OBJ }
func (a *ArrayInteger) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	for _, item := range a.Value {
		out.WriteString(fmt.Sprintf("%d", item))
		out.WriteString(",")
	}
	s := out.String()
	if s[len(s)-1] == ',' {
		s = s[0 : len(s)-1]
	}
	return s + "]"
}

// ArrayString ["a", "b", "c"]
type ArrayString struct {
	Value []string
}

func (a *ArrayString) node()                  {}
func (a *ArrayString) expressionNode()        {}
func (a *ArrayString) ObjectType() ObjectType { return ARRAY_STRING_OBJ }
func (a *ArrayString) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	for _, item := range a.Value {
		out.WriteString(item)
		out.WriteString(",")
	}
	s := out.String()
	if s[len(s)-1] == ',' {
		s = s[0 : len(s)-1]
	}
	return s + "]"
}

// CallExpression 函数调用
type CallExpression struct {
	Function  Expression   // Identifier
	Arguments []Expression //
}

func (ce *CallExpression) node()           {}
func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// PrefixExpresion 前缀表达式，<前缀运算符><表达式>
type PrefixExpresion struct {
	Operator TokenType
	Right    Expression
}

func (pe *PrefixExpresion) node()           {}
func (pe *PrefixExpresion) expressionNode() {}
func (pe *PrefixExpresion) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(string(pe.Operator))
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpresion 中缀表达式，<表达式> <中缀运算符> <表达式>
type InfixExpression struct {
	Left     Expression
	Operator TokenType
	Right    Expression
}

func (ie *InfixExpression) node()           {}
func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + string(ie.Operator) + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}
