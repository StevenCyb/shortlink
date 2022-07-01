package store

// Store is an interface for a store
type Store interface {
	Save(shortURL, originalURL string) error
	Retrieve(shortURL string) (string, error)
}
