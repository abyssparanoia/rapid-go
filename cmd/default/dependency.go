package main

import (
	"github.com/abyssparanoia/rapid-go/internal/default/handler/api"
	"github.com/abyssparanoia/rapid-go/internal/default/infrastructure/repository"
	"github.com/abyssparanoia/rapid-go/internal/default/usecase"
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluefirebaseauth"
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluemysql"
	"github.com/abyssparanoia/rapid-go/internal/pkg/httpheader"
	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/volatiletech/sqlboiler/boil"
)

// Dependency ... dependency
type Dependency struct {
	Log              *log.Middleware
	gluefirebaseauth *gluefirebaseauth.Middleware
	DummyHTTPHeader  *httpheader.Middleware
	HTTPHeader       *httpheader.Middleware
	UserHandler      *api.UserHandler
}

// Inject ... indect dependency
func (d *Dependency) Inject(e *Environment) {

	var lCli log.Writer
	var firebaseauth gluefirebaseauth.Firebaseauth

	authCli := gluefirebaseauth.NewClient(e.ProjectID)
	// fCli := gluefirestore.NewClient(e.ProjectID)

	// Config
	dbCfg := gluemysql.NewConfig()

	// pkg
	_ = gluemysql.NewClient(dbCfg)

	if e.ENV == "LOCAL" {
		lCli = log.NewWriterStdout()
		firebaseauth = gluefirebaseauth.NewDebug(authCli)
		boil.DebugMode = true
	} else {
		lCli = log.NewWriterStackdriver(e.ProjectID)
		firebaseauth = gluefirebaseauth.New(authCli)
	}

	// Repository
	uRepo := repository.NewUser()

	// Service
	dhh := httpheader.NewDummy()
	hh := httpheader.New()
	uSvc := usecase.NewUser(uRepo)

	// Middleware
	d.Log = log.NewMiddleware(lCli, e.MinLogSeverity)
	d.gluefirebaseauth = gluefirebaseauth.NewMiddleware(firebaseauth)
	d.DummyHTTPHeader = httpheader.NewMiddleware(dhh)
	d.HTTPHeader = httpheader.NewMiddleware(hh)

	// Handler
	d.UserHandler = api.NewUserHandler(uSvc)
}
