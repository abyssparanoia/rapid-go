package main

import (
	"github.com/abyssparanoia/rapid-go/src/handler/api"
	"github.com/abyssparanoia/rapid-go/src/infrastructure/repository"
	"github.com/abyssparanoia/rapid-go/src/lib/firebaseauth"
	"github.com/abyssparanoia/rapid-go/src/lib/httpheader"
	"github.com/abyssparanoia/rapid-go/src/lib/log"
	"github.com/abyssparanoia/rapid-go/src/lib/mysql"
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
	var firebaseAuthService firebaseauth.Service

	if e.ENV == "LOCAL" {
		lCli = log.NewWriterStdout()
		firebaseAuthService = firebaseauth.NewDebugService()
	} else {
		lCli = log.NewWriterStackdriver(e.ProjectID)
		firebaseAuthService = firebaseauth.NewService()
	}

	// Config
	dbCfg := mysql.NewConfig("rapid_go")

	// Lib
	dbConn := mysql.NewClient(dbCfg)

	// Repository
	uRepo := repository.NewUser(dbConn)

	// Service
	dhhSvc := httpheader.NewDummyService()
	hhSvc := httpheader.NewService()
	uSvc := service.NewUser(uRepo)

	// Middleware
	d.Log = log.NewMiddleware(lCli, e.MinLogSeverity)
	d.FirebaseAuth = firebaseauth.NewMiddleware(firebaseAuthService)
	d.DummyHTTPHeader = httpheader.NewMiddleware(dhhSvc)
	d.HTTPHeader = httpheader.NewMiddleware(hhSvc)

	// Handler
	d.UserHandler = api.NewUserHandler(uSvc)
}
