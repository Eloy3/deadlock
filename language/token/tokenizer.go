package token

import (
	"fmt"
	"strings"
)

func TokenizeProgram(content string) ([]Token, error) {
	tokens := []Token{}
	lines := strings.Split(content, "\n")

	for lineNum, line := range lines {
		lineTokens, err := TokenizeLine(line, lineNum)
		if err != nil {
			return tokens, err
		}
		tokens = append(tokens, lineTokens...)
	}

	// Append EOF token
	tokens = append(tokens, Token{Type: EOF, Literal: "EOF", Line: len(lines)})

	return tokens, nil
}

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
			tokens = append(tokens, Token{Type: SEMICOLON, Literal: string(ch), Line: lineNum})
			i++

		case ch == '=':
			if i+1 < len(line) && rune(line[i+1]) == '=' {
				tokens = append(tokens, Token{Type: EQUALS, Literal: "==", Line: lineNum})
				i += 2
			} else {
				tokens = append(tokens, Token{Type: ASSIGN, Literal: "=", Line: lineNum})
				i++
			}

		case ch == '!':
			if i+1 < len(line) && rune(line[i+1]) == '=' {
				tokens = append(tokens, Token{Type: NOT_EQUALS, Literal: "!=", Line: lineNum})
				i += 2
			} else {
				tokens = append(tokens, Token{Type: NOT, Literal: "!", Line: lineNum})
				i++
			}

		case ch == '+':
			tokens = append(tokens, Token{Type: PLUS, Literal: "+", Line: lineNum})
			i++

		case ch == '-':
			tokens = append(tokens, Token{Type: MINUS, Literal: "-", Line: lineNum})
			i++

		case ch == '*':
			tokens = append(tokens, Token{Type: MULTIPLY, Literal: "*", Line: lineNum})
			i++

		case ch == '/':
			tokens = append(tokens, Token{Type: DIVIDE, Literal: "/", Line: lineNum})
			i++

		case ch == '>':
			tokens = append(tokens, Token{Type: GREATER, Literal: ">", Line: lineNum})
			i++

		case ch == '<':
			tokens = append(tokens, Token{Type: LESS_GREATER, Literal: "<", Line: lineNum})
			i++

		case ch == '(':
			tokens = append(tokens, Token{Type: LPAREN, Literal: "(", Line: lineNum})
			i++

		case ch == ')':
			tokens = append(tokens, Token{Type: RPAREN, Literal: ")", Line: lineNum})
			i++

		case isDigit(ch):
			start := i
			for i < len(line) && isDigit(rune(line[i])) {
				i++
			}
			tokens = append(tokens, Token{Type: NUMBER, Literal: line[start:i], Line: lineNum})

		case isLetter(ch):
			id, length := readId(line[i:])
			value, exists := keywords[id]
			if exists {
				tokens = append(tokens, Token{Type: value, Literal: id, Line: lineNum})
			} else {
				tokens = append(tokens, Token{Type: IDENTIFIER, Literal: id, Line: lineNum})
			}
			i += length

		default:
			err := fmt.Errorf("token %q not recognized on line %d", ch, lineNum)
			return nil, err
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
