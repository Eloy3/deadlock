package token

const (
	ILLEGAGL   = "ILLEGAL"
	EOF        = "EOF"
	IDENTIFIER = "IDENTIFIER"
	NUMBER     = "NUMBER"
	STRING     = "STRING"

	LET    = "LET"
	SHARED = "SHARED"
	INT    = "INT"
	THREAD = "THREAD"
	LOCK   = "LOCK"
	UNLOCK = "UNLOCK"
	PRINT  = "PRINT"
	TRUE   = "true"
	FALSE  = "false"
	IF     = "IF"
	ELSE   = "ELSE"
	LOCAL  = "LOCAL"
	MUTEX  = "MUTEX"

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
	"let": LET, "true": TRUE, "false": FALSE,
	"if": IF, "else": ELSE, "local": LOCAL, "mutex": MUTEX,
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
