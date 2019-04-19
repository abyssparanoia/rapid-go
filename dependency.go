package main

import (
	"github.com/abyssparanoia/rapid-go/src/handler/api"
	"github.com/abyssparanoia/rapid-go/src/infrastructure/repository"
	"github.com/abyssparanoia/rapid-go/src/lib/firebaseauth"
	"github.com/abyssparanoia/rapid-go/src/lib/httpheader"
	"github.com/abyssparanoia/rapid-go/src/lib/mysql"
	"github.com/abyssparanoia/rapid-go/src/service"
)

// Dependency ... 依存性
type Dependency struct {
	DummyFirebaseAuth *firebaseauth.Middleware
	FirebaseAuth      *firebaseauth.Middleware
	DummyHTTPHeader   *httpheader.Middleware
	HTTPHeader        *httpheader.Middleware
	UserHandler       *api.UserHandler
}

// Inject ... 依存性を注入する
func (d *Dependency) Inject() {
	// Config
	dbCfg := mysql.GetSQLConfig()

	// Lib
	dbConn := mysql.NewSQLClient(dbCfg)

	// Repository
	uRepo := repository.NewUser(dbConn)

	// Service
	dfaSvc := firebaseauth.NewDummyService()
	faSvc := firebaseauth.NewService()
	dhhSvc := httpheader.NewDummyService()
	hhSvc := httpheader.NewService()
	uSvc := service.NewUser(uRepo)

	// Middleware
	d.DummyFirebaseAuth = firebaseauth.NewMiddleware(dfaSvc)
	d.FirebaseAuth = firebaseauth.NewMiddleware(faSvc)
	d.DummyHTTPHeader = httpheader.NewMiddleware(dhhSvc)
	d.HTTPHeader = httpheader.NewMiddleware(hhSvc)

	// Handler
	d.UserHandler = api.NewUserHandler(uSvc)
}
