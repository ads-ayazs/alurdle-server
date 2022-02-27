package store

import (
	"sync"

	"github.com/matryer/resync"
)

// WordleStore returns an object used to persist Wordle game data that implements
// the Store interface.
func WordleStore() (Store, error) {
	ws := getWordleStore()
	return ws, nil
}

// Save persists the given content object using the provided ID.
//
// Note that if the ID is already in use, it will intentionally overwrite any
// existing content with this new object. An error is returned if the ID string
// is empty or invalid.
func (s *wordleStore) Save(id string, content interface{}) error {
	// Check that the ID is valid
	if err := validateId(id); err != nil {
		return err
	}

	// Obtain write lock and store content
	s.mu.Lock()
	defer s.mu.Unlock()

	s.games[id] = content

	return nil
}

// Load returns an object previously saved using the given ID.
//
// If no content exists for the ID, then it returns nil (no error). An error is
// returned if the ID string is empty or invalid.
func (s *wordleStore) Load(id string) (interface{}, error) {
	// Validate the ID
	if err := validateId(id); err != nil {
		return nil, err
	}

	// Obtain a read lock and load the object by ID (if it exists).
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.games[id]
	if !ok {
		return nil, nil
	}

	return c, nil
}

// Exists returns true if the given ID is associated with content in the store.
//
// An error is returned if the ID string is empty or invalid.
func (s *wordleStore) Exists(id string) (bool, error) {
	// Check that the ID string is valid
	if err := validateId(id); err != nil {
		return false, err
	}

	// Obtain a read lock
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if an object exists for the given ID
	_, ok := s.games[id]
	return ok, nil
}

// Delete the object in the store for the given ID string.
//
// An error is returned if the ID string is empty or invalid.
func (s *wordleStore) Delete(id string) error {
	// Validate the ID string
	if err := validateId(id); err != nil {
		return err
	}

	// Obtain a write lock
	s.mu.Lock()
	defer s.mu.Unlock()

	// Delete the object from the store
	if _, ok := s.games[id]; ok {
		delete(s.games, id)
	} else {
		return ErrInvalidId
	}

	return nil
}

// PurgeAll deletes all objects from the store.
func (s *wordleStore) PurgeAll() error {
	// Obtain a write lock
	s.mu.Lock()
	defer s.mu.Unlock()

	// Delete all the objects
	for k, _ := range s.games {
		delete(s.games, k)
	}

	return nil
}

/////////////////

type wordleStore struct {
	mu    sync.RWMutex
	games map[string]interface{}
}

var singleStore *wordleStore
var once resync.Once // using resync.Once to facilitate testing

func getWordleStore() *wordleStore {
	if singleStore == nil {
		once.Do(
			func() {
				singleStore = new(wordleStore) //&wordleStore{}

				singleStore.mu.Lock()
				defer singleStore.mu.Unlock()

				singleStore.games = make(map[string]interface{})
			})
	}

	return singleStore
}

// Created to facilitate testing
func resetWordleStore() {
	singleStore = nil
	once.Reset()
}

func validateId(id string) error {
	if len(id) < 1 {
		return ErrInvalidId
	}

	return nil
}
