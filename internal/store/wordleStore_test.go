package store

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWordleStore(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		result *wordleStore
		err    error
	}{
		{result: &wordleStore{games: map[string]interface{}{}}, err: nil},
	}

	for _, test := range tests {
		resetWordleStore()
		s, err := WordleStore()
		assert.IsType(test.err, err, "unexpected error type")
		if err != nil {
			assert.EqualError(err, test.err.Error())
			continue // This test returned a valid error so move to the next test
		}

		if assert.NotNil(s) {
			v, ok := s.(*wordleStore)
			assert.True(ok)

			assert.NotNil(v)
			assert.NotNil(v.games)
		}
	}
}

// func (s wordleStore) Save(id string, content interface{}) error
func TestSave(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		id      string
		content string
		err     error
	}{
		{id: "1a2b3c4d5e", content: "This is the first content", err: nil},
		{id: "2a4b6c8d0e", content: "This is the second content", err: nil},
		{id: "", content: "cause an error", err: errors.New("invalid id")},
	}

	resetWordleStore() // Ensure that we get a new instance
	store, err := WordleStore()
	require.NoError(err, "error obtaining the instance")
	require.NotNil(store, "instance is nil")

	for count, test := range tests {
		err := store.Save(test.id, test.content)
		assert.IsType(test.err, err, "unexpected error type")
		if err != nil {
			assert.EqualError(err, test.err.Error())
			continue // This test returned a valid error so move to the next test
		}

		v, ok := store.(*wordleStore)
		assert.True(ok)
		assert.Equal(len(v.games), count+1)
		assert.Contains(v.games, test.id, "id is missing from map")
		// assert.Contains(v.games, test.content, "content is missing from map")
	}

}

// func (s wordleStore) Load(id string) (interface{}, error)
func TestLoad(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		id      string
		content string
		// result  wordleStore
		err error
	}{
		{id: "", content: "cause an error", err: errors.New("invalid id")},
		{id: "1a2b3c4d5e", content: "This is the first content", err: nil},
		{id: "2a4b6c8d0e", content: "This is the second content", err: nil},
	}

	resetWordleStore()
	store, err := WordleStore()
	require.NoError(err, "error obtaining the instance")
	require.NotNil(store, "instance is nil")

	// Save the test data first
	for _, test := range tests {
		if test.err != nil {
			continue
		}
		err := store.Save(test.id, test.content)
		require.NoError(err, "problem saving the test data")
	}

	// Test the Load function
	for _, test := range tests {
		content, err := store.Load(test.id)
		assert.IsType(test.err, err, "unexpected error type")
		if err != nil {
			assert.EqualError(err, test.err.Error())
			continue // This test returned a valid error so move to the next test
		}

		assert.Equal(test.content, content)
	}
}

func TestExists(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		id      string
		content string
		err     error
	}{
		{id: "", content: "cause an error", err: errors.New("invalid id")},
		{id: "1a2b3c4d5e", content: "This is the first content", err: nil},
		{id: "2a4b6c8d0e", content: "This is the second content", err: nil},
	}

	resetWordleStore()
	store, err := WordleStore()
	require.NoError(err, "error obtaining the instance")
	require.NotNil(store, "instance is nil")

	// Test for non existance
	for _, test := range tests {
		e, err := store.Exists(test.id)
		assert.IsType(test.err, err, "unexpected error type")
		if err != nil {
			assert.EqualError(err, test.err.Error())
			continue // This test returned a valid error so move to the next test
		}
		assert.False(e)
	}

	// Save the test data
	for _, test := range tests {
		if test.err != nil {
			continue
		}
		err := store.Save(test.id, test.content)
		require.NoError(err, "problem saving the test data")
	}

	// Test for existance
	for _, test := range tests {
		e, err := store.Exists(test.id)
		assert.IsType(test.err, err, "unexpected error type")
		if err != nil {
			assert.EqualError(err, test.err.Error())
			continue // This test returned a valid error so move to the next test
		}
		assert.True(e)
	}
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		id      string
		content string
		err     error
	}{
		{id: "", content: "cause an error", err: errors.New("invalid id")},
		{id: "1a2b3c4d5e", content: "This is the first content", err: nil},
		{id: "2a4b6c8d0e", content: "This is the second content", err: nil},
	}

	resetWordleStore()
	store, err := WordleStore()
	require.NoError(err, "error obtaining the instance")
	require.NotNil(store, "instance is nil")

	// Save the test data
	for _, test := range tests {
		if test.err != nil {
			continue
		}
		err := store.Save(test.id, test.content)
		require.NoError(err, "problem saving the test data")
	}

	// Test deleting items that exist
	for _, test := range tests {
		v, ok := store.(*wordleStore)
		require.True(ok)
		storeSize := len(v.games)

		err := store.Delete(test.id)
		assert.IsType(test.err, err, "unexpected error type")
		if err != nil {
			assert.EqualError(err, test.err.Error())
			assert.Equal(storeSize, len(v.games))
			continue // This test returned a valid error so move to the next test
		}

		assert.Equal(storeSize-1, len(v.games))
	}

	// Test deleting non-existant items
	for _, test := range tests {
		v, ok := store.(*wordleStore)
		require.True(ok)
		storeSize := len(v.games)

		err := store.Delete(test.id)
		assert.Error(err)
		if test.err != nil {
			assert.EqualError(err, test.err.Error())
			continue // This test returned a valid error so move to the next test
		}

		assert.ErrorIs(err, ErrInvalidId)
		assert.Equal(storeSize, len(v.games))
	}

}

func TestPurgeAll(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		id      string
		content string
		err     error
	}{
		{id: "", content: "cause an error", err: errors.New("invalid id")},
		{id: "1a2b3c4d5e", content: "This is the first content", err: nil},
		{id: "2a4b6c8d0e", content: "This is the second content", err: nil},
	}

	resetWordleStore()
	store, err := WordleStore()
	require.NoError(err, "error obtaining the instance")
	require.NotNil(store, "instance is nil")

	// Test purging empty store
	err = store.PurgeAll()
	assert.NoError(err)

	// Save the test data
	for _, test := range tests {
		if test.err != nil {
			continue
		}
		err := store.Save(test.id, test.content)
		require.NoError(err, "problem saving the test data")
	}

	// Test purging store with data
	err = store.PurgeAll()
	assert.NoError(err)

	// Ensure store is empty
	v, ok := store.(*wordleStore)
	require.True(ok)
	storeSize := len(v.games)
	assert.Zero(storeSize)

	// Test purging empty store tat has just been purged
	err = store.PurgeAll()
	assert.NoError(err)

}

// func createCleanStore() (Store, error) {
// 	store, err := WordleStore()
// 	if err != nil {
// 		return store, err
// 	}

// 	if v, ok := store.(wordleStore); ok {
// 		for k, _ := range v.games {
// 			delete(v.games, k)
// 		}
// 	}

// 	return store, err
// }
