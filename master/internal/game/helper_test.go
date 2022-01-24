package game

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateWords(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		s      string
		result string
		err    error
	}{
		{s: "", result: "", err: errors.New("invalid word length")},
		{s: "adi", result: "adi", err: errors.New("invalid word length")},
		{s: "blagu", result: "blagu", err: errors.New("word is not in dictionary")},
		{s: "blank", result: "BLANK", err: nil},
		{s: "blANk", result: "BLANK", err: nil},
		{s: "anime", result: "ANIME", err: nil},
	}

	for _, test := range tests {
		res, err := validateWord(test.s)
		if test.err != nil {
			assert.IsType(test.err, err)
			assert.EqualError(err, test.err.Error())
			assert.Equal(test.result, res, "should return original string")
			continue // This test returned a valid error so move to the next test
		}
		assert.NoError(err)

		assert.Equal(test.result, res, "return is not as expected")
	}
}
