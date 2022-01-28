package game

import (
	"strings"

	"aluance.io/wordleserver/internal/config"
	"aluance.io/wordleserver/internal/dictionary"
)

func validateWord(s string, options ...interface{}) (string, error) {
	optSecretWord := ""
	if len(options) > 0 {
		optSecretWord = strings.ToUpper(options[0].(string))
	}

	if len(s) != config.CONFIG_GAME_WORDLENGTH {
		return s, ErrWordLength
	}

	s = strings.ToUpper(s)

	// When test word and secret word are the same, no need to check the dictionary.
	if s == optSecretWord {
		return strings.ToUpper(s), nil // automatically valid
	}

	if !dictionary.IsWordValid(s) {
		return s, ErrInvalidWord
	}

	return s, nil
}
