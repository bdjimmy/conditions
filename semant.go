package conditions

import "fmt"

// semantic detection
// type check

// prefixProtos type check
var prefixProtos = map[TokenType][]ObjectType{
	BANG: {
		BOOLEAN_OBJ,
	},
}

// infixProtos type check
var infixProtos = map[TokenType]map[ObjectType]ObjectType{
	GT: {
		INTEGER_OBJ: INTEGER_OBJ,
		STRING_OBJ:  STRING_OBJ,
	},
	GT_EQUAL: {
		INTEGER_OBJ: INTEGER_OBJ,
		STRING_OBJ:  STRING_OBJ,
	},
	LT: {
		INTEGER_OBJ: INTEGER_OBJ,
		STRING_OBJ:  STRING_OBJ,
	},
	LT_EQUAL: {
		INTEGER_OBJ: INTEGER_OBJ,
		STRING_OBJ:  STRING_OBJ,
	},
	EQ: {
		INTEGER_OBJ: INTEGER_OBJ,
		STRING_OBJ:  STRING_OBJ,
		BOOLEAN_OBJ: BOOLEAN_OBJ,
	},
	NOT_EQ: {
		INTEGER_OBJ: INTEGER_OBJ,
		STRING_OBJ:  STRING_OBJ,
		BOOLEAN_OBJ: BOOLEAN_OBJ,
	},
	AND: {
		BOOLEAN_OBJ: BOOLEAN_OBJ,
	},
	OR: {
		BOOLEAN_OBJ: BOOLEAN_OBJ,
	},
}

// funcProtos type check
var funcProtos = map[string][][2][]ObjectType{
	"len": {
		{
			{STRING_OBJ},  // args
			{INTEGER_OBJ}, // return
		},
		{
			{ARRAY_STRING_OBJ}, // args
			{INTEGER_OBJ},      // return
		},
		{
			{ARRAY_INTEGER_OBJ}, // args
			{INTEGER_OBJ},       // return
		},
	},
}

func (p *Parser) CheckType(node Node) ObjectType {
	if len(p.errors) != 0 {
		return ERROR_OBJ
	}
	switch n := node.(type) {
	case *Program:
		return p.CheckType(n.Expression)
	case *Integer:
		return INTEGER_OBJ
	case *String:
		return STRING_OBJ
	case *Boolean:
		return BOOLEAN_OBJ
	case *Identifier:
		return IDENT_OBJ
	case *ArrayString:
		return ARRAY_STRING_OBJ
	case *ArrayInteger:
		return ARRAY_INTEGER_OBJ
	case *PrefixExpresion:
		{
			expects, ok := prefixProtos[n.Operator]
			if !ok {
				p.errors = append(p.errors, fmt.Sprintf("PrefixExpresion unknow operator(%s)", n.Operator))
				return ERROR_OBJ
			}
			right := p.CheckType(n.Right)

			// special case
			if right == IDENT_OBJ {
				return BOOLEAN_OBJ
			}

			for _, expect := range expects {
				if expect == right {
					return BOOLEAN_OBJ
				}
			}
			return ERROR_OBJ
		}

	case *InfixExpression:
		{
			expects, ok := infixProtos[n.Operator]
			if !ok {
				p.errors = append(p.errors, fmt.Sprintf("InfixExpression unknow operator(%s)", n.Operator))
				return ERROR_OBJ
			}
			left := p.CheckType(n.Left)
			right := p.CheckType(n.Right)

			// special case
			if left == IDENT_OBJ || right == IDENT_OBJ {
				return BOOLEAN_OBJ
			}

			rightExpect, ok := expects[left]
			if !ok {
				p.errors = append(p.errors, fmt.Sprintf("InfixExpression(%s) unknow left type(%s)",
					n.String(), left))
				return ERROR_OBJ
			}
			if rightExpect != right {
				p.errors = append(p.errors, fmt.Sprintf("InfixExpression <exp>%s<exp> right expect %s, got %s",
					n.Operator, rightExpect, right))
				return ERROR_OBJ
			}
			return BOOLEAN_OBJ
		}
	case *CallExpression:
		{
			expects, ok := funcProtos[n.Function.String()]
			if !ok {
				p.errors = append(p.errors, fmt.Sprintf("CallExpression unknow function(%s)", n.Function.String()))
				return ERROR_OBJ
			}
			if len(expects[0][0]) != len(n.Arguments) {
				p.errors = append(p.errors, fmt.Sprintf("CallExpression %s args len error, expect %d, got %d",
					n.Function.String(), len(expects), len(n.Arguments)))
				return ERROR_OBJ
			}
			returnType := ERROR_OBJ
			for _, expectArgs := range expects {
				if returnType != ERROR_OBJ {
					break
				}
				for i, expectType := range expectArgs[0] {
					actual := p.CheckType(n.Arguments[i])
					if expectType != actual && actual != IDENT_OBJ {
						break
					}
					if i == len(expectArgs[0])-1 {
						returnType = expectArgs[1][0]
					}
				}
			}
			return returnType
		}
	}
	return ERROR_OBJ
}
