package main

import (
	"github.com/abyssparanoia/rapid-go/internal/default/handler/api"
	"github.com/abyssparanoia/rapid-go/internal/default/infrastructure/repository"
	"github.com/abyssparanoia/rapid-go/internal/default/usecase"
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluefirebaseauth"
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluemysql"
	"github.com/abyssparanoia/rapid-go/internal/pkg/httpheader"
	"github.com/abyssparanoia/rapid-go/internal/pkg/httpmiddleware"
	"github.com/volatiletech/sqlboiler/boil"
	"go.uber.org/zap"
)

// Dependency ... dependency
type Dependency struct {
	httpMiddleware   *httpmiddleware.HTTPMiddleware
	gluefirebaseauth *gluefirebaseauth.Middleware
	DummyHTTPHeader  *httpheader.Middleware
	HTTPHeader       *httpheader.Middleware
	UserHandler      *api.UserHandler
}

// Inject ... indect dependency
func (d *Dependency) Inject(e *environment, logger *zap.Logger) {

	var firebaseauth gluefirebaseauth.Firebaseauth

	authCli := gluefirebaseauth.NewClient(e.ProjectID)
	// fCli := gluefirestore.NewClient(e.ProjectID)

	// pkg
	_ = gluemysql.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase)

	if e.Envrionment == "local" {
		firebaseauth = gluefirebaseauth.NewDebug(authCli)
		boil.DebugMode = true
	} else {
		firebaseauth = gluefirebaseauth.New(authCli)
	}

	// Repository
	uRepo := repository.NewUser()

	// Service
	dhh := httpheader.NewDummy()
	hh := httpheader.New()
	uSvc := usecase.NewUser(uRepo)

	// Middleware
	d.httpMiddleware = httpmiddleware.New(logger)

	d.gluefirebaseauth = gluefirebaseauth.NewMiddleware(firebaseauth)
	d.DummyHTTPHeader = httpheader.NewMiddleware(dhh)
	d.HTTPHeader = httpheader.NewMiddleware(hh)

	// Handler
	d.UserHandler = api.NewUserHandler(uSvc)
}
