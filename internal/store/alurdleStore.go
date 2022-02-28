package store

import (
	"sync"

	"github.com/matryer/resync"
)

// AlurdleStore returns an object used to persist Alurdle game data that implements
// the Store interface.
func AlurdleStore() (Store, error) {
	ws := getAlurdleStore()
	return ws, nil
}

// Save persists the given content object using the provided ID.
//
// Note that if the ID is already in use, it will intentionally overwrite any
// existing content with this new object. An error is returned if the ID string
// is empty or invalid.
func (s *alurdleStore) Save(id string, content interface{}) error {
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
func (s *alurdleStore) Load(id string) (interface{}, error) {
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
func (s *alurdleStore) Exists(id string) (bool, error) {
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
func (s *alurdleStore) Delete(id string) error {
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
func (s *alurdleStore) PurgeAll() error {
	// Obtain a write lock
	s.mu.Lock()
	defer s.mu.Unlock()

	// Delete all the objects
	for k := range s.games {
		delete(s.games, k)
	}

	return nil
}

/////////////////

// alurdleStore struct wraps the mutex and map used for in-memory persistence.
type alurdleStore struct {
	mu    sync.RWMutex
	games map[string]interface{}
}

// singleStore points at a singleton instance
var singleStore *alurdleStore
var once resync.Once // using resync.Once to facilitate testing

// getAlurdleStore is the singleton constructor for the alurdleStore.
func getAlurdleStore() *alurdleStore {
	if singleStore == nil {
		once.Do(
			func() {
				singleStore = new(alurdleStore) //&alurdleStore{}

				singleStore.mu.Lock()
				defer singleStore.mu.Unlock()

				singleStore.games = make(map[string]interface{})
			})
	}

	return singleStore
}

// resetAlurdleStore facilitates testing.
func resetAlurdleStore() {
	singleStore = nil
	once.Reset()
}

// validateId returns an error if the ID string is not valid
func validateId(id string) error {
	// Check that the string is not empty
	if len(id) < 1 {
		return ErrInvalidId
	}

	return nil
}
