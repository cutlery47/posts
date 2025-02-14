package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cutlery47/posts/config"
	"github.com/cutlery47/posts/internal/service"
	storage "github.com/cutlery47/posts/internal/storage/user-storage"
)

type authRoutes struct {
	svc *service.Service

	conf config.Handler
}

func (ar *authRoutes) handleRegister(w http.ResponseWriter, r *http.Request) {
	var (
		in storage.InUser
	)

	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("bad user data: %v", err)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := ar.svc.Register(r.Context(), in)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error when registering user: %v", err)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("internal server error: %v", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ar *authRoutes) handleLogin(w http.ResponseWriter, r *http.Request) {
	var (
		in storage.InUser
	)

	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("bad user data: %v", err)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sesh, err := ar.svc.Login(r.Context(), in)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error when loging user in: %v", err)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(sesh)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("internal server error: %v", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ar *authRoutes) handleLogout(w http.ResponseWriter, r *http.Request) {

}
