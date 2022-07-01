package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInMemory(t *testing.T) {
	expirationDuration := 1 * time.Second
	gcInterval := 1 * time.Second
	originalURL := "long_url"
	shortURL := "short_url"

	t.Run("SaveAndRetrieve_Success", func(t *testing.T) {
		store := NewInMemory(expirationDuration, gcInterval)

		store.Save(shortURL, originalURL)
		require.Len(t, store.memoryItems, 1)

		retrievedURL, err := store.Retrieve(shortURL)
		require.NoError(t, err)
		require.Equal(t, originalURL, retrievedURL)
	})

	t.Run("Expire_Success", func(t *testing.T) {
		store := NewInMemory(expirationDuration, gcInterval)

		store.Save(shortURL, originalURL)
		require.Len(t, store.memoryItems, 1)

		time.Sleep(gcInterval * 2)
		require.Len(t, store.memoryItems, 0)
	})

	t.Run("NotFound_Fail", func(t *testing.T) {
		store := NewInMemory(expirationDuration, gcInterval)

		retrievedURL, err := store.Retrieve("not_existing_url")
		require.Error(t, err)
		require.Empty(t, retrievedURL)
	})
}
