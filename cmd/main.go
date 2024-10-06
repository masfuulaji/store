package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/masfuulaji/store/internal/routes"
)

func main() {
	r := chi.NewRouter()
	routes.SetupRoutes(r)
	http.ListenAndServe(":3000", r)
}
