# 条件表达式解析器
-   针对结构体字段的校验可以转化为条件表达式的解析，最终返回true or false
-   这样的表达方式更加符合我们日常的代码习惯

```golang
package main

import (
	"fmt"
	"github.com/bdjimmy/conditions"
)

func main() {
	input := `(len(abc) > 1 && X == "123") || Y in [1, 2, 3]`
	p := conditions.NewParser(conditions.NewLexer(input))
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		// some err
	}

	env := conditions.NewEnvironment()
	env.Set("abc", &conditions.String{
		Value: "-=-=-=-",
	})
	env.Set("X", &conditions.String{
		Value: "123",
	})
	env.Set("Y", &conditions.Integer{
		Value: 3,
	})

	fmt.Println(conditions.Eval(program, env))

}
```

## 支持的数据类型
-   nil
-   int
-   string
-   boolean
-   array

## 支持的运算符
-   !<表达式>
-   <表达式> == <表达式>
-   <表达式> >  <表达式>
-   <表达式> >= <表达式>
-   <表达式> <  <表达式>
-   <表达式> <= <表达式>
-   <表达式> && <表达式>
-   <表达式> || <表达式>
-   <表达式> == <表达式>
-   

## 支持函数调用
-   len($F)
-   zero($F)
-   regexp($F, regexp string)
