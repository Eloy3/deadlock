package token

func TokenizeLine(line string, lineNum int) ([]Token, error) {
	tokens := []Token{}
	i := 0

	for i < len(line) {
		ch := rune(line[i])
		switch {
		case isWhitespace(ch):
			// Ignore whitespace
			i++

		case ch == ';':
			tokens = append(tokens, Token{Type: SEMICOLON, Literal: string(ch)})
			i++

		case isLetter(ch):
			id, length := readId(line[i:])
			value, exists := keywords[id]
			if exists {
				tokens = append(tokens, Token{Type: value, Literal: id, Line: lineNum})
			} else {
				tokens = append(tokens, Token{Type: IDENTIFIER, Literal: id, Line: lineNum})
			}
			i += length

		case symbols[ch]:
			if ch == '=' {
				tokens = append(tokens, Token{Type: EQUALS, Literal: string(ch)})
			} else {
				tokens = append(tokens, Token{Type: SYMBOL, Literal: string(ch)})
			}
			i++

		case isDigit(ch):
			start := i
			for i < len(line) && isDigit(rune(line[i])) {
				i++
			}
			tokens = append(tokens, Token{Type: NUMBER, Literal: line[start:i]})
		default:
			i++
			tokens = append(tokens, Token{Type: STRING, Literal: string(ch)})
		}

	}
	return tokens, nil
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func readId(line string) (string, int) {
	for i, ch := range line {
		if !isLetter(ch) && !isDigit(ch) {
			return line[:i], i
		}
	}
	return line, len(line)
}
