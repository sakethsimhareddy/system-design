package api

import (
	"net/http"

	"url-shortener/internal/handler"
)

func NewRouter(h *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", h.ShortenURL)
	mux.HandleFunc("/", h.Redirect)
	return mux
}
