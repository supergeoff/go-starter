package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/api", func(w http.ResponseWriter, _ *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]string{"message": "check"})
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		panic(fmt.Sprintf("Server failed to start: %v", err))
	}
}
