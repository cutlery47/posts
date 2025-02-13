package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(handler http.Handler) *chi.Mux {
	var (
		mux = chi.NewMux()
	)

	mux.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		r.Mount("/graphql", handler)
	})

	return mux
}
