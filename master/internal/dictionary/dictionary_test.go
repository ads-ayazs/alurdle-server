package dictionary

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TEST_DICTIONARY_LENGTH = 1367
const TEST_DICTIONARY_FILENAME = "google-10000-english-usa-no-swears-medium.txt"

func TestGenerateWord(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	err := Initialize(TEST_DICTIONARY_FILENAME)
	require.NoError(err)

	word, err := GenerateWord()
	assert.NoError(err)
	assert.Equal(5, len(word))
}

func TestIsWordValid(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	err := Initialize(TEST_DICTIONARY_FILENAME)
	require.NoError(err)

	// Test valid words
	testWords := []string{"anime", "bills", "drawn", "lives", "nodes"}
	for _, tw := range testWords {
		assert.True(IsWordValid(tw))
	}

	// Test invalid words
	testWords = []string{"xxxxx", "whizz", "bangs"}
	for _, tw := range testWords {
		assert.False(IsWordValid(tw))
	}
}

func TestInitialize(t *testing.T) {
	assert := assert.New(t)
	rand.Seed(time.Now().UnixNano())

	// Test standard initialization
	err := Initialize("")
	assert.NoError(err)

	// Use controled initialization
	err = Initialize(TEST_DICTIONARY_FILENAME)
	assert.NoError(err)

	// Make sure that the dictionary is not empty
	assert.NotZero(len(wordleDict.words))
	assert.Equal(TEST_DICTIONARY_LENGTH, len(wordleDict.words))
	assert.Equal(TEST_DICTIONARY_LENGTH, len(wordleDict.wordMap))

	// Test that the length of a random word is 5
	assert.Equal(5, len(wordleDict.words[rand.Intn(TEST_DICTIONARY_LENGTH)]))

	// assert.Equal(wordleDict.words[rand.Intn(TEST_DICTIONARY_LENGTH)], "bless")
}
