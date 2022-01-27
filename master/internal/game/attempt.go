package game

import (
	"bytes"
	"encoding/json"
	"time"
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

type WordleAttempt struct {
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

func (h LetterHint) String() string {
	if s, ok := mapLetterHintToString[h]; ok {
		return s
	}
	return "unknown"
}

func (h LetterHint) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	buf.WriteString(mapLetterHintToString[h])
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

func (h *LetterHint) UnmarshalJSON(b []byte) error {
	var j string
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	*h = mapStringToLetterHint[j]
	return nil
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
