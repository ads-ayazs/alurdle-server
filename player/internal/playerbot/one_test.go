package playerbot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOne(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		result oneBot
		err    error
	}{
		{result: oneBot{}, err: nil},
	}

	for _, test := range tests {
		pb, err := createOne()
		assert.NoError(err)
		assert.IsType(&test.result, pb)

		v, ok := pb.(*oneBot)
		assert.True(ok)
		assert.NotEmpty(v.id)
	}
}

func TestOneBotPlayGame(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	pb, err := createOne()
	require.NoError(err)

	ch := make(chan string)
	go pb.PlayGame(&ch)
	defer close(ch)

	select {
	case s := <-ch:
		assert.NotEmpty(s)
	case <-time.After(1000 * time.Second):
		assert.Fail("timed out without receiving from channel")
	}
}
