package api

import (
	"log"
	"net/http"
	"os"
	"shortlink/pkg/store"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Server provides the api for this service
type Server struct {
	router     *mux.Router
	store      store.Store
	httpServer *http.Server
	listen     string
}

// ListenAndServe run server (blocking)
func (s *Server) ListenAndServe(listen string) error {
	s.listen = listen
	httpServer := http.Server{
		Addr:    listen,
		Handler: handlers.CombinedLoggingHandler(os.Stdout, s.router),
	}
	s.httpServer = &httpServer

	log.Printf("Listening on %s", listen)
	return httpServer.ListenAndServe()
}

// NewServer create a new server instance
func NewServer(store store.Store) *Server {
	router := mux.NewRouter()
	server := &Server{
		router: router,
		store:  store,
	}

	router.Path("/").Methods(http.MethodGet).HandlerFunc(server.Index)
	router.Path("/").Methods(http.MethodPost).HandlerFunc(server.CreateShortUrl)
	router.Path("/{short_url}").HandlerFunc(server.ShortURLRedirect)

	return server
}
