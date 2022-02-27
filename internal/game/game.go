/*
Package game implements the Wordle game functionality.

This package functionality is intended to be exposed through a RESTful API. The
primary interface is Game.

Key functions:
	Create(secretWord) - Returns a new game, where secretWord is the five-letter
	word to be guessed.
	Retrieve(id string) - Loads and returns the game that matches the provided
	game ID.
*/
package game

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"aluance.io/wordleserver/internal/config"
	"aluance.io/wordleserver/internal/dictionary"
	"aluance.io/wordleserver/internal/store"
	"github.com/rs/xid"
)

// GameStatusEnum describes the states of a game at any time.
type GameStatusType int64

const (
	InPlay GameStatusType = iota
	Won
	Lost
	Resigned
)

// Game interface defines the basic functionality of a Game.
type Game interface {

	// Describe returns a JSON representation of the game object state (including the secret word).
	Describe() (string, error)

	// Play attempts a guess by passing in a five-letter word. Returns hints for each letter in the guess.
	Play(tryWord string) (string, error)

	// Resign ends the game before winning or losing.
	Resign() (string, error)
}

// Create is the factory function that returns a Game instance that is ready
// to play.
//
// You can either pass it the desired secret word as a parameter or pass an
// empty string to have it generate a secret word. If the secret word is
// invalid (due to length for example) or if the function fails to generate
// a word, it will return an error.
func Create(secretWord string) (Game, error) {
	// Generate a secret word if none was passed in
	if len(secretWord) < 1 {
		var err error
		if secretWord, err = dictionary.GenerateWord(); err != nil {
			return nil, err
		}
	}

	// Validate the secret word
	sw, err := validateWord(secretWord, secretWord)
	if err != nil {
		return nil, err
	}

	// Configure the game object
	game := &wordleGame{}
	game.Id = xid.New().String()
	game.SecretWord = sw
	game.Attempts = []*WordleAttempt{}
	game.Status = InPlay
	game.LastUpdated = time.Now()

	// Store the new game state
	s, err := store.WordleStore()
	if err != nil {
		return game, err
	}
	if err := s.Save(game.Id, game); err != nil {
		return game, err
	}

	return game, nil
}

// Retrieve a Game from the WordleStore by its ID.
//
// Requires a valid game ID as the parameter and will result in an error if it
// if unable to load the Game.
func Retrieve(id string) (Game, error) {
	// Load content by game ID from the WordleStore
	s, err := store.WordleStore()
	if err != nil {
		return nil, err
	}
	content, err := s.Load(id)
	if err != nil {
		return nil, err
	}

	// Cast the content to a Game object and return it
	game, ok := content.(Game)
	if !ok {
		return nil, ErrSerialization
	}

	return game, nil
}

// Describe returns a JSON representation of the game object state (including
// the secret word).
func (g *wordleGame) Describe() (string, error) {
	return g.statusReport(), nil
}

// Play attempts a guess by passing in a five-letter word and returns hints for
// each letter in the guess within the JSON game status string.
//
// Returns an error if the guessed word is not five letters in length, or if the
// game is over or has exceeded the maximum number of guess attempts.
func (g *wordleGame) Play(tryWord string) (string, error) {
	// Check the game is in play
	if g.Status != InPlay {
		return g.statusReport(), ErrGameOver
	}
	// Check that this attempt does not exeecded the maximum turns
	if len(g.Attempts) >= config.CONFIG_GAME_MAXATTEMPTS ||
		g.ValidAttempts >= config.CONFIG_GAME_MAXVALIDATTEMPTS {
		g.Status = Lost
		return g.statusReport(), ErrOutOfTurns
	}

	// Create an attempt and validate the tryWord
	attempt := g.addAttempt()
	tw, err := validateWord(tryWord, g.SecretWord)
	attempt.TryWord = tw
	if err != nil {
		attempt.IsValidWord = false

		if len(g.Attempts) >= config.CONFIG_GAME_MAXATTEMPTS ||
			g.ValidAttempts >= config.CONFIG_GAME_MAXVALIDATTEMPTS {
			g.Status = Lost
		}
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

// Resign ends the game before winning or losing.
//
// Returns the updated game status as JSON after resigning or an error if the
// resign attempt was unsuccessful (for example, if the game was already lost).
// An error can also be generated when the game resigns successfully but the
// WordleStore fails to update.
func (g *wordleGame) Resign() (string, error) {
	// TODO: Verify that the game is InPlay or return error

	// Update the game state
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

// MarshallJSON returns JSON byte representation of GameStatusType.
func (t GameStatusType) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	buf.WriteString(mapGameStatusToString[t])
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

// UnmarshallJSON hydrates JSON bytes into GameStatusType.
func (t *GameStatusType) UnmarshalJSON(b []byte) error {
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

// String returns string representation of GameStatusType
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
