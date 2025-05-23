package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/supergeoff/go-starter/apps/client/internal/handlers"
)

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	fs := http.FileServer(http.Dir("build/assets"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	r.Get("/", handlers.IndexHandler)
	return r
}

func main() {
	r := setupRouter()
	err := http.ListenAndServe(":3001", r)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
		panic("Server failed to start")
	}
}
