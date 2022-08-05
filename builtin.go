package condition

type BuiltinFunction func(args ...Object) Object
type Builtin struct {
	Fn BuiltinFunction
}

func (bf *Builtin) ObjectType() ObjectType { return FUNCTION_OBJ }
func (bf *Builtin) Inspect() string        { return "builtin function" }

// 内置函数列表
var builtins = map[string]*Builtin{}

// RegisterBuiltin registers a built-in function.  This is used to register
// our "standard library" functions.
func RegisterBuiltin(name string, fun BuiltinFunction) {
	builtins[name] = &Builtin{Fn: fun}
}

func init() {
	RegisterBuiltin("len", func(args ...Object) Object {
		if len(args) != 1 {
			return newError("wrong number of argument. got=%d, want=1", len(args))
		}
		switch arg := args[0].(type) {
		case *String:
			return &Integer{Value: int64(len(arg.Value))}
		case *ArrayString:
			return &Integer{Value: int64(len(arg.Value))}
		case *ArrayInteger:
			return &Integer{Value: int64(len(arg.Value))}
		default:
			return newError("argument to `len` not supported, got %s", args[0].ObjectType())
		}
	})
}
