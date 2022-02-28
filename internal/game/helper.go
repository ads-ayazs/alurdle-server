package game

import (
	"strings"

	"aluance.io/alurdleserver/internal/config"
	"aluance.io/alurdleserver/internal/dictionary"
)

// validateWord determines if a word is valid for a game and returns the
// checked word and an error value if it is not valid.
//
// This is a variadic function that optionally accepts the secret word as the
// second argument. If the try word matches the secret word then it is not
// checked againstthe dictionary.
func validateWord(s string, options ...interface{}) (string, error) {
	// Retrieve the optional secret word if it was passed
	optSecretWord := ""
	if len(options) > 0 {
		optSecretWord = strings.ToUpper(options[0].(string))
	}

	// Check the word length
	if len(s) != config.CONFIG_GAME_WORDLENGTH {
		return s, ErrWordLength
	}

	// Convert s to uppercase
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
