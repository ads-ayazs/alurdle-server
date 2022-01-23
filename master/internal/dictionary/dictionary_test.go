package dictionary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateWord(t *testing.T) {

}

func TestIsWordValid(t *testing.T) {

}

func TestInitialize(t *testing.T) {
	assert := assert.New(t)

	filename := "google-10000-english-usa-no-swears-medium.txt"

	err := Initialize(filename)
	assert.NoError(err)

	assert.NotZero(len(wordleDict.words))
	assert.Equal(1367, len(wordleDict.words))
	// assert.Equal(wordleDict.words[1300], "bless")
}
