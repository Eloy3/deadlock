package language

import "strings"

func Tokenize(line string) []string {
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(line, ";")
	return strings.Fields(line)
}
