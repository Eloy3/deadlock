package main

import (
	"bufio"
	"deadlock/language"
	"fmt"
	"os"
)

func main() {

	f, err := os.Open("simdata/examples/example01.dlk")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	nline := 0
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("line: ", line)
		tokens, _ := language.TokenizeLine(line, nline)
		for _, token := range tokens {
			fmt.Printf("Token: Type=%s, Value=%s\n", token.Type, token.Value)
		}
		nline++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

}
