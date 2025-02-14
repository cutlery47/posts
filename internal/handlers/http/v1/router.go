package v1

import (
	"net/http"

	"github.com/cutlery47/posts/config"
	"github.com/cutlery47/posts/internal/handlers/http/v1/auth"
	gql "github.com/cutlery47/posts/internal/handlers/http/v1/graphql"
	"github.com/cutlery47/posts/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func New(conf config.Handler, svc *service.Service) (*chi.Mux, error) {
	var (
		mux = chi.NewMux()
	)

	gql, err := gql.New(svc)
	if err != nil {
		return nil, err
	}

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Group(func(r chi.Router) {
			r.Mount("/graphql", gql)
		})

		r.Group(func(r chi.Router) {
			r.Mount("/auth", auth.New(conf, svc))
			r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
		})
	})

	return mux, nil
}
