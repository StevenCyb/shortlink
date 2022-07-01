package store

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("URL not found")
)

// memoryItem represent an item
// in the `InMemory` store
type memoryItem struct {
	shortURL        string
	originalURL     string
	expireTimestamp time.Time
}

type InMemory struct {
	lock               sync.RWMutex
	expirationDuration time.Duration
	memoryItems        []memoryItem
}

// Save a given short and original URL pair.
// This function will never return an error.
func (im *InMemory) Save(shortURL, originalURL string) error {
	im.lock.Lock()
	defer im.lock.Unlock()

	im.memoryItems = append(im.memoryItems, memoryItem{
		shortURL:        shortURL,
		originalURL:     originalURL,
		expireTimestamp: time.Now().Add(im.expirationDuration),
	})

	return nil
}

// Retrieve original URL for given short URL.
// Returns not found error if no item with given
// short URL exists.
func (im *InMemory) Retrieve(shortURL string) (string, error) {
	im.lock.RLock()
	defer im.lock.RUnlock()

	for _, memoryItem := range im.memoryItems {
		if memoryItem.shortURL == shortURL {
			return memoryItem.originalURL, nil
		}
	}

	return "", ErrNotFound
}

// NewInMemory create a new in memory store
func NewInMemory(expirationDuration, gcInterval time.Duration) *InMemory {
	store := &InMemory{
		memoryItems:        []memoryItem{},
		expirationDuration: expirationDuration,
	}

	// run the GC
	go func(im *InMemory) {
		for range time.Tick(gcInterval) {
			im.lock.Lock()

			cleanMemoryItems := []memoryItem{}
			now := time.Now()
			for _, memoryItem := range im.memoryItems {
				if memoryItem.expireTimestamp.After(now) {
					cleanMemoryItems = append(cleanMemoryItems, memoryItem)
				}
			}
			im.memoryItems = cleanMemoryItems

			im.lock.Unlock()
		}
	}(store)

	return store
}
