package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
}

type config struct {
	addr string
}

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var users = []User{{FirstName: "John", LastName: "Doe"}, {FirstName: "Jane", LastName: "Doe"}, {FirstName: "Tyler", LastName: "The Creator"}}

func (app *application) AddUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users = append(users, user)

	newMsg := fmt.Sprintf("%s %s Added!", user.FirstName, user.LastName)
	msg := struct {
		Message string `json:"msg"`
	}{Message: newMsg}

	if err := json.NewEncoder(w).Encode(msg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/v1", func(r chi.Router) {
		r.Post("/register", app.AddUser)
		r.Get("/users", app.GetUsers)
	})

	return r
}

func (app *application) run(r http.Handler) error {
	server := http.Server{
		Addr:         app.config.addr,
		Handler:      r,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("Server is running on port %s", app.config.addr)
	return server.ListenAndServe()
}
