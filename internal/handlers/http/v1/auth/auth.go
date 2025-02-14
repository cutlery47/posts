package auth

import (
	"encoding/json"
	"fmt"
	"log"
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
		log.Println(fmt.Sprintf("[REQUEST] bad user data: %v", err))
		w.Write([]byte("bad request"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sesh, err := ar.svc.Login(r.Context(), in)
	if err != nil {
		log.Println("[REQUEST] error when loging user in: ", err)
		w.Write([]byte("couldn't log you in"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(sesh)
	if err != nil {
		log.Println("[REQUEST] internal server error: ", err)
		w.Write([]byte("internal server error"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ar *authRoutes) handleLogout(w http.ResponseWriter, r *http.Request) {
	var (
		sesh storage.Session
	)

	err := json.NewDecoder(r.Body).Decode(&sesh)
	if err != nil {
		log.Println(fmt.Sprintf("[REQUEST] bad session data: %v", err))
		w.Write([]byte("bad request"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = ar.svc.Logout(r.Context(), sesh)
	if err != nil {
		log.Println("[REQUEST] error when loging user in: ", err)
		w.Write([]byte("couldn't log you in"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(sesh)
	if err != nil {
		log.Println("[REQUEST] internal server error: ", err)
		w.Write([]byte("internal server error"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
