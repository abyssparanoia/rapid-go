package main

import (
	"net/http"

	"github.com/abyssparanoia/rapid-go/src/config"
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

	// no need to authenticate for development
	if config.IsEnvDeveloping() {
		r.Route("/noauth/v1", func(r chi.Router) {
			r.Use(d.DummyFirebaseAuth.Handle)
			r.Use(d.DummyHTTPHeader.Handle)
			subRouting(r, d)
		})
	}

	// need to authenticate for production
	r.Route("/v1", func(r chi.Router) {
		r.Use(d.FirebaseAuth.Handle)
		r.Use(d.HTTPHeader.Handle)
		subRouting(r, d)
	})

	// Ping
	r.Get("/ping", handler.Ping)
	r.Get("/", handler.Ping)

	http.Handle("/", r)
}

func subRouting(r chi.Router, d Dependency) {
	// API
	r.Route("/users", func(r chi.Router) {
		//r.Post("/", d.UserHandler.Create)
		r.Get("/{userID}", d.UserHandler.Get)
	})

}
