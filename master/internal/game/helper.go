package game

import (
	"fmt"
	"strings"

	"aluance.io/wordle/master/internal/dictionary"
)

func validateWord(s string) (string, error) {
	if len(s) != 5 {
		return s, fmt.Errorf("invalid word length")
	}

	if !dictionary.IsWordValid(s) {
		return s, fmt.Errorf("word is not in dictionary")
	}

	return strings.ToUpper(s), nil
}
