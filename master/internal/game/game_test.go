package game

import (
	"encoding/json"
	"errors"
	"fmt"
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
		{secretWord: "", result: wordleGame{SecretWord: "", Attempts: []*WordleAttempt{}}, err: errors.New("invalid word length")},
		{secretWord: "adieu", result: wordleGame{SecretWord: "ADIEU", Attempts: []*WordleAttempt{}}, err: nil},
		{secretWord: "OuiJa", result: wordleGame{SecretWord: "OUIJA", Attempts: []*WordleAttempt{}}, err: nil},
	}

	for _, test := range tests {
		g, err := Create(test.secretWord)
		assert.IsType(test.err, err)
		if err != nil {
			assert.EqualError(err, test.err.Error())
			continue // This test returned a valid error so move to the next test
		}

		if assert.NotNil(g) {
			v, ok := g.(*wordleGame)
			assert.True(ok)

			assert.Equal(test.result.SecretWord, v.SecretWord, "secretWord doesn't match")
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
		{createWord: "Adieu", result: "{\"Id\":\"a0bcxxxx0x0x0x00xxxx\",\"Status\":0,\"SecretWord\":\"ADIEU\",\"Attempts\":[]}", err: nil},
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
		{createWord: "adieu", tryWord: "", result: "{}", err: errors.New("invalid word length")},
		{createWord: "adieu", tryWord: "zzzzz", result: "{\"AttemptsUsed\":1,\"GameStatus\":\"InPlay\",\"IsValidWord\":true,\"TryResult\":[3,3,3,3,3],\"TryWord\":\"ZZZZZ\"}", err: nil},
		{createWord: "adieu", tryWord: "ADIeu", result: "{\"AttemptsUsed\":1,\"GameStatus\":\"Won\",\"IsValidWord\":true,\"TryResult\":[1,1,1,1,1],\"TryWord\":\"ADIEU\",\"WinningAttempt\":1}", err: nil},
		{createWord: "adieu", tryWord: "ouija", result: "{\"AttemptsUsed\":1,\"GameStatus\":\"InPlay\",\"IsValidWord\":true,\"TryResult\":[3,2,1,3,2],\"TryWord\":\"OUIJA\"}", err: nil},
	}

	for _, test := range tests {
		game, err := Create(test.createWord)
		require.NoError(err, "Create() returned error when creating Game")
		require.NotNil(game, "unable to create a Game object")

		s, err := game.Play(test.tryWord)
		assert.IsType(test.err, err, "returned unexpected error")
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
		{createWord: "adieu", result: "{\"AttemptsUsed\":0,\"GameStatus\":\"Resigned\"}", err: nil},
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
		{createWord: "adieu", result: &WordleAttempt{TryWord: "", IsValidWord: false, TryResult: []LetterHint{0, 0, 0, 0, 0}}},
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
		{createWord: "Adieu", id: "", err: nil},
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
