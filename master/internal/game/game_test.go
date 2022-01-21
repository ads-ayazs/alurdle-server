package game

import (
	"errors"
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
		{secretWord: "", result: wordleGame{secretWord: "", attempts: nil}, err: errors.New("invalid secretWord length")},
		{secretWord: "adieu", result: wordleGame{secretWord: "adieu", attempts: nil}, err: nil},
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

			assert.Equal(test.result.secretWord, v.secretWord, "secretWord doesn't match")
			assert.Equal(test.result.attempts, v.attempts, "attempts doesn't match")
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
		{createWord: "Adieu", result: "{}", err: nil},
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

		assert.Equal(test.result, s, "returned unexpected result")
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
		{createWord: "adieu", tryWord: "", result: "{\"TryWord\":\"\",\"IsValidWord\":false,\"TryResult\":[4,4,4,4,4]}", err: nil},
	}

	for _, test := range tests {
		game, err := Create(test.createWord)
		require.NoError(err, "Create() returned error when creating Game")
		require.NotNil(game, "unable to create a Game object")

		s, err := game.Play(test.tryWord)
		assert.ErrorIs(test.err, err, "returned unexpected error")
		if err != nil {
			assert.EqualError(err, test.err.Error())
			continue // This test returned a valid error so move to the next test
		}

		assert.Equal(test.result, s, "returned unexpected result")
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
		{createWord: "adieu", result: "{}", err: nil},
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
