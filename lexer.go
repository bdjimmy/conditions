package condition

// Lexer 代表一个词法解析器
type Lexer struct {
	input        string
	position     int  // 所输入字符串中的当前位置，指向当前字符
	readPosition int  // 所输入字符串中的当前读取位置，指向当前字符的后一个字符
	ch           byte // 当前正在查看的字符
}

// New 实例化词法解析器
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}

// NextToken 从input中读取下一个token
func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{
				Type:    EQ,
				Literal: "==",
			}
		}
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case '[':
		tok = newToken(LBRACKET, l.ch)
	case ']':
		tok = newToken(RBRACKET, l.ch)
	case ';':
		tok = newToken(SEMICOLON, l.ch)
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = Token{
				Type:    AND,
				Literal: "&&",
			}
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = Token{
				Type:    OR,
				Literal: "||",
			}
		}
	case '~':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{
				Type:    REG,
				Literal: "~=",
			}
		}
	case 'i':
		if l.peekChar() == 'n' {
			l.readChar()
			tok = Token{
				Type:    IN,
				Literal: "in",
			}
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{
				Type:    NOT_EQ,
				Literal: "!=",
			}
		} else {
			tok = newToken(BANG, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{
				Type:    GT_EQUAL,
				Literal: ">=",
			}
		} else {
			tok = newToken(GT, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{
				Type:    LT_EQUAL,
				Literal: "<=",
			}
		} else {
			tok = newToken(LT, l.ch)
		}
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Type = EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) { // 数字
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal) // 关键字和用户定义的标识符区分开
			return tok
		}
		if isDigit(l.ch) { // 标识符
			tok.Type = INT
			tok.Literal = l.readNumber()
			return tok
		}
		tok = newToken(ILLEGAL, l.ch)
	}
	l.readChar()
	return tok
}

// 实例化一个token
func newToken(tokenType TokenType, ch byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// 读取input中的下一个字符，并向前移动指针
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// 读取一个标识符
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// 跳过所有的空白字符
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// 读取一个数字
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// 判断是否是一个合法的字符
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		ch == '_'
}
