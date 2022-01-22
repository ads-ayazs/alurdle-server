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
	"encoding/json"
	"fmt"
	"strings"

	"aluance.io/wordle/master/internal/store"
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
}

// Factory used to create a game
func Create(secretWord string) (Game, error) {
	sw, err := validateWord(secretWord)
	if err != nil {
		return nil, err
	}
	// game := new(wordleGame)
	game := &wordleGame{}
	game.Id = xid.New().String()
	game.SecretWord = sw
	game.Attempts = []*WordleAttempt{}
	game.Status = InPlay

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
		return nil, fmt.Errorf("content is not a game")
	}

	return game, nil
}

func (g wordleGame) Describe() (string, error) {
	gameStr := fmt.Sprint(g)
	return gameStr, nil
}

func (g *wordleGame) Play(tryWord string) (string, error) {
	if g.Status != InPlay {
		return g.statusReport(), fmt.Errorf("game is finished")
	}
	if len(g.Attempts) >= 6 {
		g.Status = Lost
		return g.statusReport(), fmt.Errorf("out of turns")
	}

	tw, err := validateWord(tryWord)
	if err != nil && len(tryWord) != 5 {
		// Not a valid attempt (word must be excatly 5 characters long)
		return "{}", err
	}

	attempt := g.addAttempt()
	attempt.TryWord = tw
	attempt.IsValidWord = true

	// Score the tryWord letters against the secret
	score := attempt.TryResult
	for i := 0; i < 5; i++ {
		idx := strings.Index(g.SecretWord, string(tw[i]))
		if idx == i {
			score[i] = Green
		} else if idx >= 0 {
			score[i] = Yellow
		} else {
			score[i] = Grey
		}
	}

	// Check for end of game conditions
	if attempt.isWinner() {
		g.Status = Won
	} else if len(g.Attempts) >= 6 {
		g.Status = Lost
	}

	// Return the attempt as JSON
	return g.turnReport(attempt), nil
}

func (g wordleGame) Resign() (string, error) {
	g.Status = Resigned

	return g.statusReport(), nil
}

/////////////

func (t GameStatusType) String() string {
	switch t {
	case InPlay:
		return "InPlay"
	case Won:
		return "Won"
	case Lost:
		return "Lost"
	case Resigned:
		return "Resigned"
	}
	return "unknown"
}

type wordleGame struct {
	Id         string
	Status     GameStatusType
	SecretWord string
	Attempts   []*WordleAttempt
}

func (g wordleGame) String() string {
	b, err := json.Marshal(g)
	if err != nil {
		return "{}"
	}

	return (string(b))
}

func (g *wordleGame) addAttempt() *WordleAttempt {
	wa := new(WordleAttempt)

	wa.TryWord = ""
	wa.IsValidWord = false
	wa.TryResult = make([]LetterHint, 5)

	g.Attempts = append(g.Attempts, wa)

	return wa
}

func (g wordleGame) statusReport() string {
	s := map[string]interface{}{}
	s["GameStatus"] = fmt.Sprint(g.Status)
	s["AttemptsUsed"] = len(g.Attempts)
	if g.Status == Won {
		s["WinningAttempt"] = len(g.Attempts)
	}

	b, err := json.Marshal(s)
	if err != nil {
		return "{}"
	}

	return string(b)
}

func (g wordleGame) turnReport(a *WordleAttempt) string {
	// Convert game status to map
	sr := g.statusReport()
	report := map[string]interface{}{}

	if err := json.Unmarshal([]byte(sr), &report); err != nil {
		return "{}"
	}

	// Convert turn status to map
	ar := fmt.Sprint(a)
	arMap := map[string]interface{}{}

	if err := json.Unmarshal([]byte(ar), &arMap); err != nil {
		return "{}"
	}

	// Comnine the maps and return as JSON
	for k, v := range arMap {
		report[k] = v
	}

	b, err := json.Marshal(report)
	if err != nil {
		return "{}"
	}

	return string(b)
}

// func (g wordleGame) copy() {}
