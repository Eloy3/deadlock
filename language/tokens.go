package language

type TokenType string

const (
	EOF        TokenType = "EOF"
	IDENTIFIER TokenType = "IDENTIFIER"
	NUMBER     TokenType = "NUMBER"
	STRING     TokenType = "STRING"
	KEYWORD    TokenType = "KEYWORD"
	OPERATOR   TokenType = "OPERATOR"
	SYMBOL     TokenType = "SYMBOL"
	COMMENT    TokenType = "COMMENT"
)

type Token struct {
	Type  TokenType
	Value string
}

var keywords = map[string]bool{
	"shared": true, "int": true, "thread": true,
	"lock": true, "unlock": true, "print": true,
}

var symbols = map[rune]bool{
	'=': true, '+': true, '-': true, '*': true, '/': true,
	';': true, '{': true, '}': true, '(': true, ')': true,
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
