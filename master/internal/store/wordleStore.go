package store

import (
	"sync"

	"github.com/matryer/resync"
)

func WordleStore() (Store, error) {
	ws := getWordleStore()
	return ws, nil
}

func (s *wordleStore) Save(id string, content interface{}) error {
	if err := validateId(id); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.games[id] = content

	return nil
}

func (s *wordleStore) Load(id string) (interface{}, error) {
	if err := validateId(id); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.games[id]
	if !ok {
		return nil, nil
	}

	return c, nil
}

func (s *wordleStore) Exists(id string) (bool, error) {
	if err := validateId(id); err != nil {
		return false, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.games[id]
	return ok, nil
}

func (s *wordleStore) Delete(id string) error {
	if err := validateId(id); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.games[id]; ok {
		delete(s.games, id)
	} else {
		return ErrInvalidId
	}

	return nil
}

func (s *wordleStore) PurgeAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()

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
