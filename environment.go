package conditions

import (
	"fmt"
	"os"
)

type Environment struct {
	store    map[string]Object
	readOnly map[string]struct{}
}

func NewEnvironment() *Environment {
	return &Environment{
		store:    make(map[string]Object),
		readOnly: map[string]struct{}{},
	}
}

func (env *Environment) Get(name string) (Object, bool) {
	obj, ok := env.store[name]
	return obj, ok
}

func (env *Environment) Set(name string, val Object) Object {
	if _, ok := env.readOnly[name]; ok {
		fmt.Printf("Attempting to modify '%s' denied; it was defined as a constant.\n", name)
		os.Exit(3)
	}
	env.store[name] = val
	return val
}

func (env *Environment) SetReadOnly(name string, val Object) Object {
	env.store[name] = val
	env.readOnly[name] = struct{}{}
	return val
}
