package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		secretWord string
		result     wordleGame
		err        error
	}{
		{secretWord: "", result: wordleGame{SecretWord: "", Attempts: []*WordleAttempt{}}, err: nil},
		{secretWord: "happy", result: wordleGame{SecretWord: "HAPPY", Attempts: []*WordleAttempt{}}, err: nil},
		{secretWord: "hAPpY", result: wordleGame{SecretWord: "HAPPY", Attempts: []*WordleAttempt{}}, err: nil},
	}

	for _, test := range tests {
		g, err := Create(test.secretWord)
		assert.IsType(test.err, err, fmt.Sprintf("\"%s\": %s", test.secretWord, err))
		if err != nil {
			assert.EqualError(err, test.err.Error())
			continue // This test returned a valid error so move to the next test
		}

		if assert.NotNil(g) {
			v, ok := g.(*wordleGame)
			assert.True(ok)

			if len(test.result.SecretWord) > 0 {
				assert.Equal(test.result.SecretWord, v.SecretWord, "secretWord doesn't match")
			}
			assert.Equal(test.result.Attempts, v.Attempts, "attempts doesn't match")
		}
	}
}

func TestDescribe(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		createWord string
		result     string
		err        error
	}{
		{createWord: "Happy", result: "{\"Id\":\"a0bcxxxx0x0x0x00xxxx\",\"Status\":0,\"SecretWord\":\"HAPPY\",\"Attempts\":[]}", err: nil},
	}

	for _, test := range tests {
		game, err := Create(test.createWord)
		require.NoError(err, "Create() returned error when creating Game")
		require.NotNil(game, "unable to create a Game object")

		s, err := game.Describe()
		assert.ErrorIs(test.err, err, "returned unexpected error")
		if err != nil {
			assert.EqualError(err, test.err.Error())
			continue // This test returned a valid error so move to the next test
		}

		sMap := map[string]interface{}{}
		require.NoError(json.Unmarshal([]byte(s), &sMap))
		for k, v := range sMap {
			assert.Contains(test.result, k)
			if k == "Id" {
				continue
			}
			assert.Contains(test.result, fmt.Sprint(v))
		}
	}
}

func TestPlay(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		createWord string
		tryWord    string
		result     string
		err        error
	}{
		{createWord: "happy", tryWord: "", result: "{}", err: errors.New("invalid word length")},
		{createWord: "happy", tryWord: "zzzzz", result: "{\"AttemptsUsed\":0,\"GameStatus\":\"InPlay\",\"IsValidWord\":false,\"TryResult\":[3,3,3,3,3],\"TryWord\":\"ZZZZZ\"}", err: errors.New("word is not in dictionary")},
		{createWord: "happy", tryWord: "happy", result: "{\"AttemptsUsed\":1,\"GameStatus\":\"Won\",\"IsValidWord\":true,\"TryResult\":[1,1,1,1,1],\"TryWord\":\"HAPPY\",\"WinningAttempt\":1}", err: nil},
		{createWord: "happy", tryWord: "puppy", result: "{\"AttemptsUsed\":1,\"GameStatus\":\"InPlay\",\"IsValidWord\":true,\"TryResult\":[3,3,1,1,1],\"TryWord\":\"PUPPY\"}", err: nil},
		{createWord: "happy", tryWord: "bless", result: "{\"AttemptsUsed\":1,\"GameStatus\":\"InPlay\",\"IsValidWord\":true,\"TryResult\":[3,3,3,3,3],\"TryWord\":\"BLESS\"}", err: nil},
	}

	for _, test := range tests {
		game, err := Create(test.createWord)
		require.NoError(err, "Create() returned error when creating Game")
		require.NotNil(game, "unable to create a Game object")

		s, err := game.Play(test.tryWord)
		assert.IsType(test.err, err, fmt.Sprintf("\"%s\" unexpected error: %s", test.tryWord, err))
		if err != nil {
			assert.EqualError(err, test.err.Error(), "returned unexpected error")
			continue // This test returned a valid error so move to the next test
		}

		assert.Equal(test.result, s, "returned unexpected result")

		v, ok := game.(*wordleGame)
		require.True(ok)
		assert.NotZero(len(v.Attempts))
	}
}

func TestResign(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		createWord string
		result     string
		err        error
	}{
		{createWord: "bless", result: "{\"AttemptsUsed\":0,\"GameStatus\":\"Resigned\"}", err: nil},
	}

	for _, test := range tests {
		game, err := Create(test.createWord)
		require.NoError(err, "Create() returned error when creating Game")
		require.NotNil(game, "unable to create a Game object")

		s, err := game.Resign()
		assert.ErrorIs(test.err, err, "returned unexpected error")
		if err != nil {
			assert.EqualError(err, test.err.Error())
			continue // This test returned a valid error so move to the next test
		}

		assert.Equal(test.result, s, "returned unexpected result")
	}
}

//func (g *wordleGame) addAttempt() *WordleAttempt
func TestAddAttempt(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		createWord string
		result     *WordleAttempt
	}{
		{createWord: "proxy", result: &WordleAttempt{TryWord: "", IsValidWord: false, TryResult: []LetterHint{0, 0, 0, 0, 0}}},
	}
	for _, test := range tests {
		game, err := Create(test.createWord)
		require.NoError(err, "Create() returned error when creating Game")
		require.NotNil(game, "unable to create a Game object")

		v, ok := game.(*wordleGame)
		require.True(ok)
		res := v.addAttempt()
		assert.NotZero(len(v.Attempts))
		assert.Equal(res.TryWord, test.result.TryWord)
		assert.Equal(res.IsValidWord, test.result.IsValidWord)
		assert.Equal(res.TryResult, test.result.TryResult)
	}
}

//func Retrieve(id string) (*Game, error)
func TestRetrieve(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		createWord string
		id         string
		err        error
	}{
		{createWord: "seven", id: "", err: nil},
	}

	// Create test games
	for i, test := range tests {
		game, err := Create(test.createWord)
		require.NoError(err, "Create() returned error when creating Game")
		require.NotNil(game, "unable to create a Game object")

		v, ok := game.(*wordleGame)
		require.True(ok)
		tests[i].id = v.Id
	}

	// Test retrieving games
	for _, test := range tests {
		game, err := Retrieve(test.id)
		assert.IsType(test.err, err, "returned unexpected error")
		if err != nil {
			assert.EqualError(err, test.err.Error())
			continue // This test returned a valid error so move to the next test
		}

		assert.NoError(err)
		assert.NotNil(game)

		v, ok := game.(*wordleGame)
		require.True(ok)
		assert.Equal(test.id, v.Id)
	}
}

func TestScoreWord(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		createWord string
		tryWord    string
		result     []LetterHint
		err        error
	}{
		{createWord: "happy", tryWord: "seven", result: []LetterHint{Grey, Grey, Grey, Grey, Grey}, err: nil},
		{createWord: "happy", tryWord: "heave", result: []LetterHint{Green, Grey, Yellow, Grey, Grey}, err: nil},
		{createWord: "happy", tryWord: "paint", result: []LetterHint{Yellow, Green, Grey, Grey, Grey}, err: nil},
		{createWord: "happy", tryWord: "peppy", result: []LetterHint{Grey, Grey, Green, Green, Green}, err: nil},
		{createWord: "happy", tryWord: "happy", result: []LetterHint{Green, Green, Green, Green, Green}, err: nil},
		{createWord: "knoll", tryWord: "wooly", result: []LetterHint{Grey, Grey, Green, Green, Grey}, err: nil},
	}

	for _, test := range tests {
		// Create test game
		game, err := Create(test.createWord)
		require.NoError(err, "Create() returned error when creating Game")
		require.NotNil(game, "unable to create a Game object")

		v, ok := game.(*wordleGame)
		require.True(ok)

		r := []LetterHint{Blank, Blank, Blank, Blank, Blank}
		err = v.scoreWord(strings.ToUpper(test.tryWord), &r)
		if test.err != nil {
			assert.Error(err)
			assert.EqualError(test.err, err.Error())
		}
		assert.NoError(err, fmt.Sprintf("\"%s\": %s", test.tryWord, err))
		assert.Exactly(test.result, r, fmt.Sprintf("\"%s\": score results do not match", test.tryWord))
	}

}
