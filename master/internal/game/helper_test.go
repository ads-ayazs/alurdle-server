package game

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

//func validateWord(s string) (string, error)
func TestValidateWords(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		s      string
		result string
		err    error
	}{
		{s: "", result: "", err: errors.New("invalid word length")},
		{s: "adi", result: "adi", err: errors.New("invalid word length")},
		{s: "ouiJa", result: "OUIJA", err: nil},
	}

	for _, test := range tests {
		res, err := validateWord(test.s)
		assert.IsType(test.err, err)
		if err != nil {
			assert.EqualError(err, test.err.Error())
			assert.Equal(test.result, res, "should return original string")
			continue // This test returned a valid error so move to the next test
		}

		assert.Equal(test.result, res, "return is not as expected")
	}
}
