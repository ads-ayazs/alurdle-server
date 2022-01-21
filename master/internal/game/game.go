package game

import (
	"encoding/json"
	"fmt"
)

// Enum for letter hints returned after each attempt
type LetterHint int

const (
	Blank  LetterHint = iota // not set
	Green                    // correct letter and position
	Yellow                   // correct letter, wrong position
	Grey                     // incorrect letter
	Red                      // invalid word
)

// Game interface
type Game interface {
	Describe() (string, error)
	Play(tryWord string) (string, error)
	Resign() (string, error)
}

// Factory used to create a game
func Create(secretWord string) (Game, error) {
	if len(secretWord) != 5 {
		return nil, fmt.Errorf("invalid secretWord length")
	}

	game := new(wordleGame)
	game.secretWord = secretWord

	return game, nil
}

func (h LetterHint) String() string {
	switch h {
	case Blank:
		return "Blank"
	case Green:
		return "Green"
	case Yellow:
		return "Yellow"
	case Grey:
		return "Grey"
	case Red:
		return "Red"
	}
	return "unknown"
}

type WordleAttempt struct {
	TryWord     string
	IsValidWord bool
	TryResult   []LetterHint
}

func (a *WordleAttempt) String() string {
	b, err := json.Marshal(a)
	if err != nil {
		return "{}"
	}

	return (string(b))
}

type wordleGame struct {
	secretWord string
	attempts   []WordleAttempt
}

func (g wordleGame) Describe() (string, error) {
	return "{}", nil
}

func (g wordleGame) Play(tryWord string) (string, error) {
	attempt := newWordleAttempt()
	attempt.TryWord = tryWord

	// Word length must be 5 otherwise it is invalid
	if len(tryWord) != 5 {
		attempt.IsValidWord = false
		for i := 0; i < 5; i++ {
			attempt.TryResult[i] = Red
		}
	} else {
		// TODO Check if tryWord is valid in dictionary
		attempt.IsValidWord = true
	}

	attemptStr := fmt.Sprint(attempt)
	return attemptStr, nil
}

func (g wordleGame) Resign() (string, error) {
	return "{}", nil
}

func newWordleAttempt() *WordleAttempt {
	wa := new(WordleAttempt)

	wa.TryWord = ""
	wa.IsValidWord = false
	wa.TryResult = make([]LetterHint, 5)

	return wa
}
