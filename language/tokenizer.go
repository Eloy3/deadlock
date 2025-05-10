package language

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
			tokens = append(tokens, Token{Type: SYMBOL, Value: string(ch)})
			i++

		case isLetter(ch):
			id, length := readId(line[i:])
			if keywords[id] {
				tokens = append(tokens, Token{Type: KEYWORD, Value: id})
			} else {
				tokens = append(tokens, Token{Type: IDENTIFIER, Value: id})
			}
			i += length

		case symbols[ch]:
			tokens = append(tokens, Token{Type: SYMBOL, Value: string(ch)})
			i++

		default:
			i++
			tokens = append(tokens, Token{Type: STRING, Value: string(ch)})
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

// int x
func readId(line string) (string, int) {
	for i, ch := range line {
		if !isLetter(ch) && !isDigit(ch) {
			return line[:i], i
		}
	}
	return line, len(line)

}
