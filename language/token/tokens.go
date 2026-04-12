package token

const (
	EOF        = "EOF"
	IDENTIFIER = "IDENTIFIER"
	NUMBER     = "NUMBER"
	STRING     = "STRING"

	SHARED = "SHARED"
	INT    = "INT"
	THREAD = "THREAD"
	LOCK   = "LOCK"
	UNLOCK = "UNLOCK"
	PRINT  = "PRINT"

	ASSIGN = "="

	EQUALS          = "=="
	NOT_EQUALS      = "!="
	GREATER         = ">"
	GREATER_EQ      = ">="
	LESS_GREATER    = "<"
	LESS_GREATER_EQ = "<="

	PLUS     = "+"
	MINUS    = "-"
	DIVIDE   = "/"
	MULTIPLY = "*"
	MODULE   = "%"

	SYMBOL  = "SYMBOL"
	COMMENT = "COMMENT"

	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	DOT       = "."

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	OR  = "||"
	AND = "&&"
	NOT = "!"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

var keywords = map[string]TokenType{
	"shared": SHARED, "int": INT, "thread": THREAD,
	"lock": LOCK, "unlock": UNLOCK, "print": PRINT,
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
