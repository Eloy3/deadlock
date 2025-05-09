package main

import (
	"bufio"
	"deadlock/language"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hulloooo")

	f, err := os.Open("simdata/examples/example01.dlk")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()
	// Read the file and process it
	// For example, you can read the file line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := language.Tokenize(line)
		fmt.Println("line: ", tokens, "\n")
		for i, token := range tokens {
			fmt.Printf("Token %d: %s\n", i, token)
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	// Example of using the Tokenize function

}
