package store

import (
	"fmt"
	"sync"
)

func WordleStore() (Store, error) {
	ws := getWordleStore()
	return ws, nil
}

func (s wordleStore) Save(id string, content interface{}) error {
	if err := validateId(id); err != nil {
		return err
	}
	s.games[id] = content

	return nil
}

func (s wordleStore) Load(id string) (interface{}, error) {
	if err := validateId(id); err != nil {
		return nil, err
	}

	c, ok := s.games[id]
	if !ok {
		return nil, nil
	}

	return c, nil
}

func (s wordleStore) Exists(id string) (bool, error) {
	if err := validateId(id); err != nil {
		return false, err
	}

	_, ok := s.games[id]
	return ok, nil
}

func (s wordleStore) Delete(id string) error {
	if err := validateId(id); err != nil {
		return err
	}

	if _, ok := s.games[id]; ok {
		delete(s.games, id)
	} else {
		return fmt.Errorf("id does not exist")
	}

	return nil
}

func (s wordleStore) PurgeAll() error {
	for k, _ := range s.games {
		delete(s.games, k)
	}

	return nil
}

/////////////////

type wordleStore struct {
	games map[string]interface{}
}

var singleStore *wordleStore
var once sync.Once

func getWordleStore() *wordleStore {
	if singleStore == nil {
		once.Do(
			func() {
				singleStore = new(wordleStore) //&wordleStore{}
				singleStore.games = make(map[string]interface{})
			})
	}

	return singleStore
}

func validateId(id string) error {
	if len(id) < 1 {
		return fmt.Errorf("invalid id")
	}

	return nil
}
