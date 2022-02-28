package game

import (
	"bytes"
	"encoding/json"
	"time"
)

// LetterHint enum for letter hints returned after each attempt.
type LetterHint int

const (
	Blank  LetterHint = iota // not set
	Green                    // correct letter and position
	Yellow                   // correct letter, wrong position
	Grey                     // incorrect letter
	Red                      // invalid word
)

// AlurdleAttempt holds the state of an attempt (guess) and its result.
type AlurdleAttempt struct {
	TryWord     string       `json:"tryWord"`
	IsValidWord bool         `json:"isValidWord"`
	TryResult   []LetterHint `json:"tryResult"`
	TimeStamp   time.Time    `json:"timeStamp"`
}

var mapLetterHintToString = map[LetterHint]string{
	Blank:  "Blank",
	Green:  "Green",
	Yellow: "Yellow",
	Grey:   "Grey",
	Red:    "Red",
}

var mapStringToLetterHint = map[string]LetterHint{
	"Blank":  Blank,
	"Green":  Green,
	"Yellow": Yellow,
	"Grey":   Grey,
	"Red":    Red,
}

// String converts LetterHint into its string representation
func (h LetterHint) String() string {
	if s, ok := mapLetterHintToString[h]; ok {
		return s
	}
	return "unknown"
}

// MarshallJSON returns JSON byte representation of a LetterHint.
func (h LetterHint) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	buf.WriteString(mapLetterHintToString[h])
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshallJSON hydrates JSON bytes into a LetterHint.
func (h *LetterHint) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*h = mapStringToLetterHint[s]
	return nil
}

// isWinner returns true if this is a winning AlurdleAttempt.
func (a AlurdleAttempt) isWinner() bool {
	winner := true

	for _, r := range a.TryResult {
		if r != Green {
			winner = false
			break
		}
	}

	return winner
}
