package conditions

import "fmt"

func Eval(node Node, env *Environment) Object {
	switch node := node.(type) {
	case *Program:
		return evalProgram(node, env)
	case *Integer:
		return node
	case *String:
		return node
	case *ArrayString:
		return node
	case *ArrayInteger:
		return node
	case *Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *Identifier:
		return evalIdentifier(node, env)
	case *CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		ret := applyFunction(function, args)
		return ret
	case *PrefixExpresion:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixOperatorExpression(node.Operator, right)
	case *InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	}
	return nil
}

func evalProgram(program *Program, env *Environment) Object {
	return Eval(program.Expression, env)
}

func evalIdentifier(node Node, env *Environment) Object {
	ident := node.(*Identifier)
	if val, ok := env.Get(ident.Value); ok {
		return val
	}
	if builtin, ok := builtins[ident.Value]; ok {
		return builtin
	}
	return newError("identifier not found: " + ident.Value)
}

// 执行前缀表达式
func evalPrefixOperatorExpression(operator TokenType, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	default:
		// 错误处理
		return newError("unknow operator:%s%s", operator, right.ObjectType())
	}
}

// 执行 !<expression>
func evalBangOperatorExpression(right Object) Object {
	switch right {
	case boolTrue:
		return boolFalse
	case boolFalse:
		return boolTrue
	}
	return boolFalse
}

func applyFunction(fn Object, args []Object) Object {
	switch fn := fn.(type) {
	case *Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.ObjectType())
	}
}

func evalExpressions(exps []Expression, env *Environment) []Object {
	var result []Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalInfixExpression(operator TokenType, left, right Object) Object {
	switch {
	case left.ObjectType() == INTEGER_OBJ && right.ObjectType() == INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.ObjectType() == STRING_OBJ && right.ObjectType() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == AND:
		return nativeBoolToBooleanObject(objectToNativeBoolean(left) && objectToNativeBoolean(right))
	case operator == OR:
		return nativeBoolToBooleanObject(objectToNativeBoolean(left) || objectToNativeBoolean(right))
	case operator == IN:
		return evalINInfixExpress(left, right)
	default:
		return newError("unknow operator: %s %s %s",
			left.ObjectType(), operator, right.ObjectType())
	}
}

func evalINInfixExpress(left, right Object) Object {
	switch {
	case left.ObjectType() == INTEGER_OBJ && right.ObjectType() == ARRAY_INTEGER_OBJ:
		leftVal := left.(*Integer).Value
		rightVal := right.(*ArrayInteger).Value
		for _, val := range rightVal {
			if leftVal == val {
				return &Boolean{
					Value: true,
				}
			}
		}
		return &Boolean{
			Value: false,
		}
	case left.ObjectType() == STRING_OBJ && right.ObjectType() == ARRAY_STRING_OBJ:
		leftVal := left.(*String).Value
		rightVal := right.(*ArrayString).Value
		for _, val := range rightVal {
			if leftVal == val {
				return &Boolean{
					Value: true,
				}
			}
		}
		return &Boolean{
			Value: false,
		}
	default:
		return newError("unknow operator: %s %s %s",
			left.ObjectType(), "IN", right.ObjectType())
	}

}

func evalIntegerInfixExpression(operator TokenType, left, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value
	switch operator {
	case "+":
		return &Integer{
			Value: leftVal + rightVal,
		}
	case "-":
		return &Integer{
			Value: leftVal - rightVal,
		}
	case "*":
		return &Integer{
			Value: leftVal * rightVal,
		}
	case "/":
		return &Integer{
			Value: leftVal / rightVal,
		}
	case "<":
		return &Boolean{
			Value: leftVal < rightVal,
		}
	case ">":
		return &Boolean{
			Value: leftVal > rightVal,
		}
	case "==":
		return &Boolean{
			Value: leftVal == rightVal,
		}
	case "!=":
		return &Boolean{
			Value: leftVal != rightVal,
		}
	default:
		return newError("unknow operator: %s %s %s",
			left.ObjectType(), operator, right.ObjectType())
	}
}

func evalStringInfixExpression(operator TokenType, left, right Object) Object {
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value
	switch operator {
	case "+":
		return &String{
			Value: leftVal + rightVal,
		}
	case "<":
		return &Boolean{
			Value: len(leftVal) < len(rightVal),
		}
	case ">":
		return &Boolean{
			Value: len(leftVal) > len(rightVal),
		}
	case "==":
		return &Boolean{
			Value: leftVal == rightVal,
		}
	case "!=":
		return &Boolean{
			Value: leftVal != rightVal,
		}
	default:
		return newError("unknow operator: %s %s %s",
			left.ObjectType(), operator, right.ObjectType())
	}
}

var (
	boolTrue  = &Boolean{Value: true}
	boolFalse = &Boolean{Value: false}
)

func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return boolTrue
	}
	return boolFalse
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.ObjectType() == ERROR_OBJ
	}
	return false
}

func newError(format string, args ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(format, args...),
	}
}

// convert object to boolean
func objectToNativeBoolean(o Object) bool {
	switch obj := o.(type) {
	case *Boolean:
		return obj.Value
	case *String:
		return obj.Value != ""
	case *Integer:
		if obj.Value == 0 {
			return false
		}
		return true
	default:
		return true
	}
}
