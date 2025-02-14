package auth

import (
	"github.com/cutlery47/posts/config"
	"github.com/cutlery47/posts/internal/service"
	"github.com/go-chi/chi/v5"
)

func New(conf config.Handler, svc *service.Service) *chi.Mux {
	var (
		mux = chi.NewMux()
	)

	auth := &authRoutes{
		svc:  svc,
		conf: conf,
	}

	mux.Group(func(r chi.Router) {
		r.Get("/register", auth.handleRegister)
		r.Get("/login", auth.handleLogin)
		r.Get("/logout", auth.handleLogout)
	})

	return mux
}
