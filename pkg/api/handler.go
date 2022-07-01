package api

import (
	"fmt"
	"net/http"
	"shortlink/pkg/shortener"
	"shortlink/pkg/store"

	"github.com/gorilla/mux"
)

// @Summary       Create short URL.
// @Description   Create a new short URL for given URL.
// @Tags          URL Shorten API
// @Success       200 {string} string "HTML code"
// @Router        / [get]
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	page := `
		<form action="/" method="post" name="form">
			<p><label for="long_url"> URL:</label>
			<input type="text" name="long_url" id="long_url" placeholder="www...." required></p>

			<input value="Submit" type="submit">
		</form>
	`

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(page))
}

// @Summary       Index page.
// @Description   Index page to enter URL.
// @Tags          Webpage
// @Param         spec body UrlCreationRequest true "The URL to shorten"
// @Success       200 {string} string "Your short URl is <a href='...'>...</a>"
// @Failure       400 {string} string "400 - invalid input format"
// @Failure       500 {string} string "500 - failed to shorten URL"
// @Failure       500 {string} string "500 - failed to store shorten URL"
// @Router        / [post]
func (s *Server) CreateShortUrl(w http.ResponseWriter, r *http.Request) {
	r.PostFormValue("long_uri")
	longURL := r.PostForm.Get("long_url")
	if longURL == "" {
		http.Error(w, "400 - invalid input format", http.StatusBadRequest)
		return
	}

	shortURL, err := shortener.Shorten(longURL)
	if err != nil {
		http.Error(w, "500 - failed to shorten URL", http.StatusInternalServerError)
		return
	}

	err = s.store.Save(shortURL, longURL)
	if err != nil {
		http.Error(w, "500 - failed to store shorten URL", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(
		fmt.Sprintf("Your short URl is <a href=\"/%s\">%s/%s</a>",
			shortURL, s.listen, shortURL),
	))
}

// @Summary       Redirect using short URL.
// @Description   Redirect using short URL.
// @Tags          URL Shorten API
// @Param         short_url path string true "Short URL postfix"
// @Header        302 {string} Location  "/..."
// @Failure       404 {string} string "404 - given short URL not found or expired"
// @Failure       500 {string} string "500 - failed to retrieve URL"
// @Router        /{short_url} [get]
func (s *Server) ShortURLRedirect(w http.ResponseWriter, r *http.Request) {
	originalURL, err := s.store.Retrieve(mux.Vars(r)["short_url"])
	if err == store.ErrNotFound {
		http.Error(w, "404 - given short URL not found or expired", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "500 - failed to retrieve URL", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, 302)
}
