package repl

import (
	"strings"
)


func CleanInput(text string) []string {
	text = strings.ToLower(text)
	text, _ = strings.CutPrefix(text, " ")
	text, _ = strings.CutSuffix(text, " ")
	words := strings.Split(text, " ")
	return words
}
