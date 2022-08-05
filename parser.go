package conditions

import (
	"fmt"
	"strconv"
)

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

// 来标记运算符的优先级
const (
	_           int = iota
	LOWEST          // 最低优先级
	COND            // AND OR
	ASSIGNS         // =
	EQUALS          // == or !=
	LESSGREATER     // > or >=  or <  or <=
	SUM             //	+ or -
	PRODUCT         // * or /
	PREFIX          // - !
	CALL            // Function(X)

)

// 运算符优先级列表, 用于将词法单元类型与其优先级相关联
var precedences = map[TokenType]int{
	AND:      COND,        // &&
	OR:       COND,        // ||
	EQ:       EQUALS,      // ==
	NOT_EQ:   EQUALS,      // !=
	LT:       LESSGREATER, // <
	LT_EQUAL: LESSGREATER, // <=
	GT:       LESSGREATER, // >
	GT_EQUAL: LESSGREATER, // >=
	IN:       PRODUCT,     // IN
	LPAREN:   CALL,        // ()
}

// Parser 递归下降语法分析器
type Parser struct {
	l              *Lexer                      // 词法分析器的实例
	curToken       Token                       // 当前正在检测的词法单元，决定下一步该做什么
	peekToken      Token                       // 下一个需要检测的词法单元
	errors         []string                    // 记录语法解析过程中的错误
	prefixParseFns map[TokenType]prefixParseFn // 前缀表达式处理函数
	infixParseFns  map[TokenType]infixParseFn  // 中缀表达式处理函数
}

// NewParser
func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []string{},
		prefixParseFns: make(map[TokenType]prefixParseFn),
		infixParseFns:  make(map[TokenType]infixParseFn),
	}

	// 注册表达式解析函数, 前缀运算符
	p.registerPrefix(IDENT, p.parseIdentifier)         // abc
	p.registerPrefix(INT, p.parseInteger)              // 123
	p.registerPrefix(STRING, p.parseString)            // "abc"
	p.registerPrefix(TRUE, p.parseBoolean)             // true
	p.registerPrefix(FALSE, p.parseBoolean)            // false
	p.registerPrefix(LBRACKET, p.parseArray)           // [
	p.registerPrefix(BANG, p.presePrefixExpression)    // !
	p.registerPrefix(LPAREN, p.parseGroupedExpression) // (

	// 注册表达式解析函数, 中缀运算符
	p.registerInfix(EQ, p.parseInfixExpression)       // ==
	p.registerInfix(NOT_EQ, p.parseInfixExpression)   // !=
	p.registerInfix(LT, p.parseInfixExpression)       // <
	p.registerInfix(LT_EQUAL, p.parseInfixExpression) // <=
	p.registerInfix(GT, p.parseInfixExpression)       // >
	p.registerInfix(GT_EQUAL, p.parseInfixExpression) // >=
	p.registerInfix(IN, p.parseInfixExpression)       // IN
	p.registerInfix(AND, p.parseInfixExpression)      // AND
	p.registerInfix(OR, p.parseInfixExpression)       // OR
	p.registerInfix(LPAREN, p.parseCallExpression)    // fn(a,b,c)

	// 读取两个词法单元，用来设置curToken和peekToken
	p.nextToken()
	p.nextToken()
	return p
}

// ParseProgram
func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	for p.curToken.Type != EOF {
		stmt := p.parseExpression(LOWEST)
		if stmt != nil {
			program.Expression = stmt
		}
		p.nextToken()
	}
	// 进行类型检测
	p.CheckType(program)

	return program
}

// 解析表达式
func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	// 向下递归
	for !p.peekTokenIs(EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// 查看当前操作符的优先级
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) presePrefixExpression() Expression {
	expression := &PrefixExpresion{
		Operator: TokenType(p.curToken.Literal),
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Operator: TokenType(p.curToken.Literal),
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Value: p.curToken.Literal}
}

// 解析一个整形的字面量
func (p *Parser) parseInteger() Expression {
	lit := &Integer{}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

// 解析字符串字面量
func (p *Parser) parseString() Expression {
	return &String{Value: p.curToken.Literal}
}

// 解析bool类型的字面量
func (p *Parser) parseBoolean() Expression {
	return &Boolean{Value: p.curTokenIs(TRUE)}
}

func (p *Parser) parseArray() Expression {
	// empty array
	if p.peekTokenIs(RBRACKET) {
		p.nextToken()
		p.errors = append(p.errors, "empty array is not allowed")
		return nil
	}
	// array type
	switch {
	case p.peekTokenIs(STRING):
		arr := &ArrayString{
			Value: make([]string, 0),
		}
		for !p.peekTokenIs(RBRACKET) {
			p.nextToken()
			arr.Value = append(arr.Value, p.curToken.Literal)
			if p.peekTokenIs(COMMA) {
				p.nextToken()
			} else {
				p.expectPeek(RBRACKET)
				break
			}
		}
		return arr
	case p.peekTokenIs(INT):
		arr := &ArrayInteger{
			Value: make([]int64, 0),
		}
		for !p.peekTokenIs(RBRACKET) {
			p.nextToken()
			i, e := strconv.ParseInt(p.curToken.Literal, 10, 64)
			if e != nil {
				p.errors = append(p.errors,
					fmt.Sprintf("the data type of the array is not a integer, err(%s)", e))
				return nil
			}
			arr.Value = append(arr.Value, i)
			if p.peekTokenIs(COMMA) {
				p.nextToken()
			} else {
				p.expectPeek(RBRACKET)
				break
			}
		}
		return arr
	default:
		p.errors = append(p.errors, "unknow array type")
		return nil
	}
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []Expression {
	args := []Expression{}

	if p.peekTokenIs(RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	p.errors = append(p.errors,
		fmt.Sprintf("no prefix parse function for %s found", t))
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

// Errors  message during syntax parsing
func (p *Parser) Errors() []string {
	return p.errors
}
