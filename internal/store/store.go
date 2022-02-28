/*
Package store provides thread-safe content persistence.

The Store interface is intended for safe use by multiple concurrent threads to
store, retreieve and delete content with unique client-generated IDs.

WordleStore is a thread-safe implementation of Store for wordle game data.

Key functions:
	WordleStore() - Obtain a Store object to manage content persistence for Wordle
	game data.
*/
package store

// Store interface describes functionality for thread-safe object persistence.
type Store interface {
	// Save persists the given content object using the provided ID.
	Save(id string, content interface{}) error

	// Load returns an object previously saved using the given ID.
	Load(id string) (interface{}, error)

	// Exists returns true if the given ID is associated with content in the store.
	Exists(id string) (bool, error)

	// Delete the object in the store for the given ID string.
	Delete(id string) error

	// PurgeAll deletes all objects from the store.
	PurgeAll() error
}
