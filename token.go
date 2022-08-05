package conditions

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// special mark
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// identifier + literal
	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"
	STRING TokenType = "STRING"

	// operator
	BANG     TokenType = "!"
	LT       TokenType = "<"
	LT_EQUAL TokenType = "<="
	GT       TokenType = ">"
	GT_EQUAL TokenType = ">="
	EQ       TokenType = "=="
	NOT_EQ   TokenType = "!="
	REG      TokenType = "~="
	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"
	IN       TokenType = "in"
	AND      TokenType = "&&"
	OR       TokenType = "||"

	// break
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	// keyword
	TRUE  TokenType = "TRUE"
	FALSE TokenType = "FALSE"
)

var keywords = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
	"in":    IN,
}

// reserved keyword
var reserved = map[string]struct{}{
	"break":       {},
	"case":        {},
	"chan":        {},
	"const":       {},
	"continue":    {},
	"default":     {},
	"defer":       {},
	"else":        {},
	"fallthrough": {},
	"for":         {},
	"func":        {},
	"go":          {},
	"goto":        {},
	"if":          {},
	"import":      {},
	"interface":   {},
	"map":         {},
	"package":     {},
	"range":       {},
	"return":      {},
	"select":      {},
	"struct":      {},
	"switch":      {},
	"type":        {},
	"var":         {},
	// special string
	"F": {},
}

// LookupIdent
func LookupIdent(ident string) TokenType {
	if t, ok := keywords[ident]; ok {
		return t
	}
	if _, ok := reserved[ident]; ok {
		return ILLEGAL
	}
	return IDENT
}
