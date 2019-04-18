package main

import (
	"fmt"
	"net/http"

	"github.com/abyssparanoia/rapid-go/src/dependency"
	"github.com/abyssparanoia/rapid-go/src/lib/log"
	"github.com/abyssparanoia/rapid-go/src/routing"
	"github.com/go-chi/chi"
)

var addr = ":3001"

func main() {
	// Dependency
	d := &dependency.Dependency{}
	d.Inject()

	// Routing
	r := chi.NewRouter()
	routing.Routing(r, d)

	// Run

	fmt.Printf("[START] server. port: %s\n", addr)
	http.ListenAndServe(addr, log.Logger(r))
}
