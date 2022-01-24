package game

import (
	"fmt"
	"strings"

	"aluance.io/wordle/internal/config"
	"aluance.io/wordle/internal/dictionary"
)

func validateWord(s string, options ...interface{}) (string, error) {
	optSecretWord := ""
	if len(options) > 0 {
		optSecretWord = strings.ToUpper(options[0].(string))
	}

	if len(s) != config.CONFIG_GAME_WORDLENGTH {
		return s, fmt.Errorf("invalid word length")
	}

	s = strings.ToUpper(s)

	// When test word and secret word are the same, no need to check the dictionary.
	if s == optSecretWord {
		return strings.ToUpper(s), nil // automatically valid
	}

	if !dictionary.IsWordValid(s) {
		return s, fmt.Errorf("word is not in dictionary")
	}

	return s, nil
}
