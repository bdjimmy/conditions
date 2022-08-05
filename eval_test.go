package conditions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	input := `len(abc) > 1 && S == 123`
	p := NewParser(NewLexer(input))
	program := p.ParseProgram()
	p.CheckType(program)

	fmt.Println(p.errors)
	assert.Equal(t, 0, len(p.errors))

	env := NewEnvironment()
	env.Set("abc", &ArrayInteger{
		Value: []int64{1, 2, 3, 4, 5, 6},
	})
	env.Set("S", &Integer{
		Value: 123,
	})
	fmt.Println(Eval(program, env))
}
