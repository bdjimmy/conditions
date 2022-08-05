# 条件表达式
-   针对结构体字段的校验可以转化为条件表达式的解析，最终返回true or false
-   这样的表达方式更加符合我们日常的代码习惯

```golang
    type RequestParams struct {
        ID int          `validate:"$ in [1, 2, 3]"`                // ID字段的取值必须是1 or 2 or 3
        Appname string  `validate:"len($) > 10 && len($) < 20"`   // Appname字段的取值的长度必须大于10并且小于20
    }
```

## 支持的数据类型
-   nil
-   int
-   float
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
