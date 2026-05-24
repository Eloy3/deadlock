package main

import (
	"bufio"
	"fmt"
	"os"

	"deadlock/language/parser"
	"deadlock/language/semantic"
	"deadlock/language/token"
)

func main() {

	dir, _ := os.Getwd()
	fmt.Println("Working directory:", dir)

	f, err := os.Open("/home/karthala/repos/deadlock/simdata/examples/example02.dlk")
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
	program := parser.ParseProgram()

	fmt.Println("\n=== AST Tree ===")
	program.PrintTree()

	// Perform semantic analysis
	fmt.Println("\n=== Semantic Analysis ===")
	symTable, errList := semantic.AnalyzeProgram(&program)

	if errList.HasErrors() {
		fmt.Println("Semantic errors found:")
		for i, err := range errList {
			fmt.Printf("  %d. %s\n", i+1, err)
		}
		return
	}

	fmt.Println("Semantic analysis passed!")
	fmt.Println("\n=== Symbol Table ===")
	fmt.Println(symTable.String())
}
