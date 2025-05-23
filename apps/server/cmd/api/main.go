package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/supergeoff/go-starter/apps/server/internal/handler"
)

// setupRouter configures and returns the chi router.
func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api", handler.ApiHandler) // handler.ApiHandler is already tested separately
	return r
}

func main() {
	r := setupRouter()
	slog.Info("Server starting on :3000") // Added a log message
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
		panic("Server failed to start")
	}
}
