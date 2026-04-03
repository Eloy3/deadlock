package main

import (
	"bufio"
	"fmt"
	"os"

	"deadlock/language/parser"
	"deadlock/language/token"
)

func main() {

	dir, _ := os.Getwd()
	fmt.Println("Working directory:", dir)

	f, err := os.Open("../simdata/examples/example01.dlk")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	nline := 0
	tokens := []token.Token{}
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("line: ", line)
		lineTokens, _ := token.TokenizeLine(line, nline)
		for _, token := range lineTokens {
			tokens = append(tokens, token)
			fmt.Printf("Token: Type=%s, Value=%s\n", token.Type, token.Literal)
		}
		nline++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	parser := parser.NewParser(tokens)
	parser.ParseProgram()
}
