package dictionary

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"aluance.io/wordle/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TEST_DICTIONARY_LENGTH = 1367
const TEST_DICTIONARY_FILENAME = "google-10000-english-usa-no-swears-medium.txt"
const TEST_DICTIONARY_FILEPATH = "data/google-10000-english-usa-no-swears-medium.txt"

func TestGenerateWord(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	wordleDict.reset()
	err := Initialize(TEST_DICTIONARY_FILEPATH)
	require.NoError(err)

	word, err := GenerateWord()
	assert.NoError(err)
	assert.Equal(config.CONFIG_GAME_WORDLENGTH, len(word))
}

func TestIsWordValid(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	wordleDict.reset()
	err := Initialize(TEST_DICTIONARY_FILEPATH)
	require.NoError(err)

	// Test valid words
	testWords := []string{"blank", "blANk", "anime", "drawn", "lives", "nodes"}
	for _, tw := range testWords {
		assert.True(IsWordValid(tw), fmt.Sprintf("\"%s\" should be valid", tw))
	}

	// Test invalid words
	testWords = []string{"xxxxx", "whizz", "bangs", "blagu"}
	for _, tw := range testWords {
		assert.False(IsWordValid(tw), fmt.Sprintf("\"%s\" should NOT be valid", tw))
	}
}

func TestInitialize(t *testing.T) {
	assert := assert.New(t)
	rand.Seed(time.Now().UnixNano())

	// Test standard initialization
	wordleDict.reset()
	assert.False(wordleDict.initalized)
	err := Initialize("")
	assert.NoError(err)

	// Make sure that the dictionary is not empty
	assert.True(wordleDict.initalized)
	assert.NotZero(len(wordleDict.words))

	// Use controled initialization
	wordleDict.reset()
	assert.False(wordleDict.initalized)
	err = Initialize(TEST_DICTIONARY_FILEPATH)
	assert.NoError(err)

	// Make sure that the dictionary is not empty
	assert.True(wordleDict.initalized)
	assert.NotZero(len(wordleDict.words))
	assert.Equal(TEST_DICTIONARY_LENGTH, len(wordleDict.words))
	assert.Equal(TEST_DICTIONARY_LENGTH, len(wordleDict.wordMap))

	// Test the length of a random word
	assert.Equal(config.CONFIG_GAME_WORDLENGTH, len(wordleDict.words[rand.Intn(TEST_DICTIONARY_LENGTH)]))

	// assert.Equal(wordleDict.words[rand.Intn(TEST_DICTIONARY_LENGTH)], "bless")
}
