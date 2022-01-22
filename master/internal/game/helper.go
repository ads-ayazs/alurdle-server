package game

import (
	"fmt"
	"strings"
)

func validateWord(s string) (string, error) {
	if len(s) != 5 {
		return s, fmt.Errorf("invalid word length")
	}

	// TODO Check if tryWord is valid in dictionary

	return strings.ToUpper(s), nil
}
