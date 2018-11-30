package main

import (
	"github.com/abyssparanoia/gke-beego/api/src/handler/api"
	"github.com/abyssparanoia/gke-beego/api/src/lib/firebaseauth"
	"github.com/abyssparanoia/gke-beego/api/src/lib/httpheader"
	"github.com/abyssparanoia/gke-beego/api/src/repository"
	"github.com/abyssparanoia/gke-beego/api/src/service"
)

// Dependency ... 依存性
type Dependency struct {
	DummyFirebaseAuth *firebaseauth.Middleware
	FirebaseAuth      *firebaseauth.Middleware
	DummyHTTPHeader   *httpheader.Middleware
	HTTPHeader        *httpheader.Middleware
	SampleHandler     *api.SampleHandler
}

// Inject ... 依存性を注入する
func (d *Dependency) Inject() {
	// Config
	// dbCfg := config.GetCSQLConfig("sample")

	// Lib
	// dbConn := cloudsql.NewCSQLClient(dbCfg)

	// Repository
	repo := repository.NewSample(nil)

	// Service
	dfaSvc := firebaseauth.NewDummyService()
	faSvc := firebaseauth.NewService()
	dhhSvc := httpheader.NewDummyService()
	hhSvc := httpheader.NewService()
	svc := service.NewSample(repo)

	// Middleware
	d.DummyFirebaseAuth = firebaseauth.NewMiddleware(dfaSvc)
	d.FirebaseAuth = firebaseauth.NewMiddleware(faSvc)
	d.DummyHTTPHeader = httpheader.NewMiddleware(dhhSvc)
	d.HTTPHeader = httpheader.NewMiddleware(hhSvc)

	// Handler
	d.SampleHandler = api.NewSampleHandler(svc)
}
