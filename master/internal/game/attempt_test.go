package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalJSON(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		lh     LetterHint
		result string
		err    error
	}{
		{lh: 100, result: "\"\"", err: nil},
		{lh: Blank, result: "\"Blank\"", err: nil},
		{lh: Green, result: "\"Green\"", err: nil},
		{lh: Yellow, result: "\"Yellow\"", err: nil},
		{lh: Grey, result: "\"Grey\"", err: nil},
		{lh: Red, result: "\"Red\"", err: nil},
	}

	for _, test := range tests {
		b, err := test.lh.MarshalJSON()
		if test.err != nil {
			assert.ErrorIs(err, test.err)
		}
		assert.NoError(err)
		assert.Equal(test.result, string(b))
	}
}

// func (h *LetterHint) UnmarshalJSON(b []byte) error {
func TestUnMarshalJSON(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		lh     string
		result LetterHint
		err    error
	}{
		{lh: "\"\"", result: Blank, err: nil},
		{lh: "\"Blank\"", result: Blank, err: nil},
		{lh: "\"Green\"", result: Green, err: nil},
		{lh: "\"Yellow\"", result: Yellow, err: nil},
		{lh: "\"Grey\"", result: Grey, err: nil},
		{lh: "\"Red\"", result: Red, err: nil},
	}

	for _, test := range tests {
		var lh LetterHint
		err := lh.UnmarshalJSON([]byte(test.lh))
		if test.err != nil {
			assert.ErrorIs(err, test.err)
		}
		assert.NoError(err)
		assert.Equal(test.result, lh)
	}
}
