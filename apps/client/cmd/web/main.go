package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/supergeoff/go-starter/apps/client/internal/pages"
)

func main() {
	r := chi.NewRouter()
	fs := http.FileServer(http.Dir("build/assets"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		pages.Home(w, r)
	})
	http.ListenAndServe(":3001", r)
}
