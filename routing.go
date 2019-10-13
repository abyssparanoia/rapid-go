package main

import (
	"net/http"

	"github.com/abyssparanoia/rapid-go/src/handler"
	"github.com/abyssparanoia/rapid-go/src/middleware"
	"github.com/go-chi/chi"
)

// Routing ... define routing
func Routing(r chi.Router, d Dependency) {

	// access control
	r.Use(middleware.AccessControl)

	// request log
	r.Use(d.Log.Handle)

	// need to authenticate for production
	r.Route("/v1", func(r chi.Router) {
		r.With(d.FirebaseAuth.Handle).Route("/users", func(r chi.Router) {
			//r.Post("/", d.UserHandler.Create)
			r.Get("/{userID}", d.UserHandler.Get)
		})
	})

	// Ping
	r.Get("/ping", handler.Ping)
	r.Get("/", handler.Ping)

	http.Handle("/", r)
}
