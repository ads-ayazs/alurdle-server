/*
Package game implements the Wordle game functionality.

This package is intended to be exposed through a RESTful API.

The primary interface is Game.

Key functions:
	Create(secretWord) - Returns a new game, where secretWord is the five-letter word to be guessed.

	Game.Play(tryWord)	- Attempt a guess by passing in a five-letter word. Returns hints for each letter in the guess.
	Game.Resign() - End the game before winning or losing.
	Game.Describe() - Returns a represantation of the game object state (including the secret word).

*/

package game

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"aluance.io/wordle/internal/config"
	"aluance.io/wordle/internal/dictionary"
	"aluance.io/wordle/internal/store"
	"github.com/rs/xid"
)

// Game status enum
type GameStatusType int64

const (
	InPlay GameStatusType = iota
	Won
	Lost
	Resigned
)

// Game interface
type Game interface {
	Describe() (string, error)
	Play(tryWord string) (string, error)
	Resign() (string, error)
	// State() (string, error)
}

// Factory used to create a game
func Create(secretWord string) (Game, error) {
	if len(secretWord) < 1 {
		var err error
		if secretWord, err = dictionary.GenerateWord(); err != nil {
			return nil, err
		}
	}

	sw, err := validateWord(secretWord, secretWord)
	if err != nil {
		return nil, err
	}
	game := &wordleGame{}
	game.Id = xid.New().String()
	game.SecretWord = sw
	game.Attempts = []*WordleAttempt{}
	game.Status = InPlay
	game.LastUpdated = time.Now()

	s, err := store.WordleStore()
	if err != nil {
		return game, err
	}
	if err := s.Save(game.Id, game); err != nil {
		return game, err
	}

	return game, nil
}

func Retrieve(id string) (Game, error) {
	s, err := store.WordleStore()
	if err != nil {
		return nil, err
	}
	content, err := s.Load(id)
	if err != nil {
		return nil, err
	}

	game, ok := content.(Game)
	if !ok {
		return nil, ErrSerialization
	}

	return game, nil
}

func (g wordleGame) Describe() (string, error) {
	return g.statusReport(), nil
}

func (g *wordleGame) Play(tryWord string) (string, error) {
	if g.Status != InPlay {
		return g.statusReport(), ErrGameOver
	}
	if len(g.Attempts) >= config.CONFIG_GAME_MAXATTEMPTS ||
		g.ValidAttempts >= config.CONFIG_GAME_MAXVALIDATTEMPTS {
		g.Status = Lost
		return g.statusReport(), ErrOutOfTurns
	}

	attempt := g.addAttempt()
	tw, err := validateWord(tryWord, g.SecretWord)
	attempt.TryWord = tw
	if err != nil {
		attempt.IsValidWord = false
		if err == ErrWordLength {
			return g.statusReport(), err
		}
		return g.statusReport(), err
	}
	attempt.IsValidWord = true
	g.ValidAttempts++

	// Score the tryWord letters against the secret
	score := attempt.TryResult
	if err := g.scoreWord(tw, &score); err != nil {
		return g.statusReport(), err
	}

	// Check for end of game conditions
	if attempt.isWinner() {
		g.Status = Won
	} else if len(g.Attempts) >= config.CONFIG_GAME_MAXATTEMPTS ||
		g.ValidAttempts >= config.CONFIG_GAME_MAXVALIDATTEMPTS {
		g.Status = Lost
	}

	g.LastUpdated = time.Now()

	// Save to game store
	gs, err := store.WordleStore()
	if err != nil {
		return g.statusReport(), err
	}
	err = gs.Save(g.Id, g)
	if err != nil {
		return g.statusReport(), err
	}

	// Return the attempt as JSON
	return g.statusReport(), nil
}

func (g *wordleGame) Resign() (string, error) {
	g.Status = Resigned
	g.LastUpdated = time.Now()

	// Save to game store
	gs, err := store.WordleStore()
	if err != nil {
		return g.statusReport(), err
	}
	err = gs.Save(g.Id, g)
	if err != nil {
		return g.statusReport(), err
	}

	return g.statusReport(), nil
}

func (t GameStatusType) Marshall() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	buf.WriteString(mapGameStatusToString[t])
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

func (t *GameStatusType) Unmarshall(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*t = mapStringToGameStatus[s]
	return nil
}

/////////////

var mapGameStatusToString = map[GameStatusType]string{
	InPlay:   "InPlay",
	Won:      "Won",
	Lost:     "Lost",
	Resigned: "Resigned",
}

var mapStringToGameStatus = map[string]GameStatusType{
	"InPlay":   InPlay,
	"Won":      Won,
	"Lost":     Lost,
	"Resigned": Resigned,
}

func (t GameStatusType) String() string {
	if s, ok := mapGameStatusToString[t]; ok {
		return s
	}
	return "unknown"
}

type wordleGame struct {
	Id            string           `json:"id"`
	Status        GameStatusType   `json:"gameStatus"`
	SecretWord    string           `json:"secretWord"`
	Attempts      []*WordleAttempt `json:"attempts"`
	ValidAttempts int              `json:"validAttempts"`
	LastUpdated   time.Time        `json:"lastUpdated"`
}

func (g *wordleGame) addAttempt() *WordleAttempt {
	wa := new(WordleAttempt)

	wa.TimeStamp = time.Now()
	wa.TryWord = ""
	wa.IsValidWord = false
	wa.TryResult = make([]LetterHint, config.CONFIG_GAME_WORDLENGTH)

	g.Attempts = append(g.Attempts, wa)
	g.LastUpdated = time.Now()

	return wa
}

func (g wordleGame) statusReport() string {
	b, err := json.Marshal(g)
	if err != nil {
		return "{}"
	}

	s := map[string]interface{}{}
	err = json.Unmarshal(b, &s)
	if err != nil {
		return "{}"
	}

	s["attemptsUsed"] = len(g.Attempts)
	if g.Status == InPlay {
		delete(s, "secretWord")
	}
	if g.Status == Won {
		s["winningAttempt"] = len(g.Attempts)
	}

	b, err = json.Marshal(s)
	if err != nil {
		return "{}"
	}

	return string(b)
}

func (g wordleGame) scoreWord(tryWord string, result *[]LetterHint) error {
	if result == nil {
		return ErrNilResult
	}
	score := *result

	// Rules for scoring:
	// 1. If the correct letter is in the correct location, mark it green
	// 2. If the letter is correct but in an incorrect location, mark it
	//    yellow UNLESS the same letter is also provided in the correct location.
	// 3. No letter should be marked yellow or green more times than it occurs
	//    in the secret word.
	// 4. Remaining unmarked letters must be marked grey.
	//
	for i := 0; i < config.CONFIG_GAME_WORDLENGTH; i++ {
		if g.SecretWord[i] == byte(tryWord[i]) {
			score[i] = Green // exact match
			continue
		} else if count := strings.Count(g.SecretWord, string(tryWord[i])); count > 0 {
			// Letter is definitely in the secret word. Check if there are other instances of the
			// same letter that are or will be marked green or yellow elsewhere in the word.
			if countLeft := strings.Count(g.SecretWord[:i+1], string(tryWord[i])); countLeft > 0 {
				// If letter occured fewer times in tryWord than secret, mark is yellow
				if strings.Count(tryWord[:i+1], string(tryWord[i])) <= countLeft {
					score[i] = Yellow
					continue
				}
			}
			if countRight := strings.Count(g.SecretWord[i:], string(tryWord[i])); countRight > 0 {
				if strings.Count(tryWord[i:], string(tryWord[i])) <= countRight {
					score[i] = Yellow
					continue
				}
			}
		}
		score[i] = Grey
	}

	return nil
}
