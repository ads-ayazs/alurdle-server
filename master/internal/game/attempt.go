package game

import "encoding/json"

// Enum for letter hints returned after each attempt
type LetterHint int

const (
	Blank  LetterHint = iota // not set
	Green                    // correct letter and position
	Yellow                   // correct letter, wrong position
	Grey                     // incorrect letter
	Red                      // invalid word
)

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

func (a WordleAttempt) isWinner() bool {
	winner := true

	for _, r := range a.TryResult {
		if r != Green {
			winner = false
			break
		}
	}

	return winner
}
