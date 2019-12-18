package main

import (
	"github.com/abyssparanoia/rapid-go/src/handler/api"
	"github.com/abyssparanoia/rapid-go/src/infrastructure/repository"
	"github.com/abyssparanoia/rapid-go/src/pkg/firebaseauth"
	"github.com/abyssparanoia/rapid-go/src/pkg/httpheader"
	"github.com/abyssparanoia/rapid-go/src/pkg/log"
	"github.com/abyssparanoia/rapid-go/src/pkg/mysql"
	"github.com/abyssparanoia/rapid-go/src/service"
)

// Dependency ... dependency
type Dependency struct {
	Log             *log.Middleware
	FirebaseAuth    *firebaseauth.Middleware
	DummyHTTPHeader *httpheader.Middleware
	HTTPHeader      *httpheader.Middleware
	UserHandler     *api.UserHandler
}

// Inject ... indect dependency
func (d *Dependency) Inject(e *Environment) {

	var lCli log.Writer
	var firebaseAuth firebaseauth.Firebaseauth

	if e.ENV == "LOCAL" {
		lCli = log.NewWriterStdout()
		firebaseAuth = firebaseauth.NewDebug()
	} else {
		lCli = log.NewWriterStackdriver(e.ProjectID)
		firebaseAuth = firebaseauth.New()
	}

	// Config
	dbCfg := mysql.NewConfig()

	// pkg
	dbConn := mysql.NewClient(dbCfg)

	// Repository
	uRepo := repository.NewUser(dbConn)

	// Service
	dhh := httpheader.NewDummy()
	hh := httpheader.New()
	uSvc := service.NewUser(uRepo)

	// Middleware
	d.Log = log.NewMiddleware(lCli, e.MinLogSeverity)
	d.FirebaseAuth = firebaseauth.NewMiddleware(firebaseAuth)
	d.DummyHTTPHeader = httpheader.NewMiddleware(dhh)
	d.HTTPHeader = httpheader.NewMiddleware(hh)

	// Handler
	d.UserHandler = api.NewUserHandler(uSvc)
}
