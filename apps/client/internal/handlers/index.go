package handlers

import (
	"net/http"

	"github.com/supergeoff/go-starter/apps/client/internal/pages"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	pages.Home(w, r)
}
