package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"shortlink/pkg/store"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	duration := 2 * time.Second
	store := store.NewInMemory(duration, duration)
	server := NewServer(store)
	server.listen = "localhost:8080"

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	server.Index(recorder, req)
	require.Equal(t, 200, recorder.Code)
	require.Equal(t, []string([]string{"text/html; charset=utf-8"}),
		recorder.HeaderMap["Content-Type"])
}

func TestCreateShortUrl(t *testing.T) {
	duration := 2 * time.Second
	store := store.NewInMemory(duration, duration)
	server := NewServer(store)
	server.listen = "localhost:8080"

	t.Run("ShortenLink_Success", func(t *testing.T) {
		longUrl := "https://github.com/StevenCyb"

		req, err := http.NewRequest(http.MethodPost, "", nil)
		require.NoError(t, err)

		form := url.Values{}
		form.Add("long_url", longUrl)
		req.PostForm = form

		recorder := httptest.NewRecorder()
		server.CreateShortUrl(recorder, req)
		require.Equal(t, 200, recorder.Code)

		require.Equal(t, []string([]string{"text/html; charset=utf-8"}),
			recorder.HeaderMap["Content-Type"])
		body, err := ioutil.ReadAll(recorder.Body)
		require.NoError(t, err)
		require.Equal(t, "Your short URl is <a href=\"/X5jF2EZJ\">localhost:8080/X5jF2EZJ</a>", string(body))
	})

	t.Run("InvalidBodyFormat_Fail", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "", nil)
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		server.CreateShortUrl(recorder, req)
		require.Equal(t, 400, recorder.Code)
	})
}

func TestShortURLRedirect(t *testing.T) {
	duration := 2 * time.Second
	store := store.NewInMemory(duration, duration)
	server := NewServer(store)
	server.listen = "localhost:8080"
	store.Save("X5jF2EZJ", "https://github.com/StevenCyb")

	req, err := http.NewRequest(http.MethodGet, "/X5jF2EZJ", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"short_url": "X5jF2EZJ",
	})

	recorder := httptest.NewRecorder()
	server.ShortURLRedirect(recorder, req)
	require.Equal(t, 302, recorder.Code)
	require.Equal(t, []string([]string{"text/html; charset=utf-8"}),
		recorder.HeaderMap["Content-Type"])
}
