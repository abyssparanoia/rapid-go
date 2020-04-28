package main

import (
	"net/http"

	"github.com/abyssparanoia/rapid-go/internal/pkg/accesscontrol"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/handler"
	"github.com/go-chi/chi"
)

// Routing ... define routing
func Routing(r chi.Router, d Dependency) {

	// access control
	r.Use(accesscontrol.Handle)

	// request log
	r.Use(d.httpMiddleware.Handle)

	// need to authenticate for production
	r.Route("/v1", func(r chi.Router) {
		r.Route("/tokens", func(r chi.Router) {
			r.Post("/set", d.TokenHandler.Set)
			r.Post("/delete", d.TokenHandler.Delete)
		})

		r.Route("/messages", func(r chi.Router) {
			r.Post("/send_to_user", d.MessageHandler.SendToUser)
			r.Post("/send_to_multi_user", d.MessageHandler.SendToMultiUser)
			r.Post("/send_to_all_user", d.MessageHandler.SendToAllUser)
		})
	})

	// Ping
	r.Get("/ping", handler.Ping)
	r.Get("/", handler.Ping)

	http.Handle("/", r)
}
