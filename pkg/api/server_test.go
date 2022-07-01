package api

import (
	"context"
	"net/http"
	"shortlink/pkg/store"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	server := NewServer(store.NewInMemory(0, 0))

	go func(server *Server) {
		err := server.ListenAndServe(":8888")
		require.Equal(t, http.ErrServerClosed, err)
	}(server)

	time.Sleep(3 * time.Second)
	server.httpServer.Shutdown(context.Background())
}
