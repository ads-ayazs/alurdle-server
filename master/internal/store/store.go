package store

type Store interface {
	Save(id string, content interface{}) error
	Load(id string) (interface{}, error)
	Exists(id string) (bool, error)
	Delete(id string) error
	PurgeAll() error
}
