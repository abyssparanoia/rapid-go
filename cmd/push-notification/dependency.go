package main

import (
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluefcm"
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluefirestore"
	"github.com/abyssparanoia/rapid-go/internal/pkg/httpheader"
	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/handler/api"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/infrastructure/repository"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/usecase"
)

// Dependency ... dependency
type Dependency struct {
	Log             *log.Middleware
	DummyHTTPHeader *httpheader.Middleware
	HTTPHeader      *httpheader.Middleware
	TokenHandler    *api.TokenHandler
	MessageHandler  *api.MessageHandler
}

// Inject ... inject dependency
func (d *Dependency) Inject(e *Environment) {
	var lCli log.Writer

	if e.ENV == "LOCAL" {
		lCli = log.NewWriterStdout()
	} else {
		lCli = log.NewWriterStackdriver(e.ProjectID)
	}

	fcmClient := gluefcm.NewClient(e.ProjectID)
	firestoreClient := gluefirestore.NewClient(e.ProjectID)

	tokenRepository := repository.NewToken(firestoreClient)
	fcmRepository := repository.NewFcm(fcmClient, e.FcmServerKey)

	tokenService := service.NewToken(tokenRepository)

	tokenUsecase := usecase.NewToken(fcmRepository, tokenRepository, tokenService)
	messageUsecase := usecase.NewMessage(fcmRepository, tokenRepository)

	dhh := httpheader.NewDummy()
	hh := httpheader.New()

	d.Log = log.NewMiddleware(lCli, e.MinLogSeverity)
	d.DummyHTTPHeader = httpheader.NewMiddleware(dhh)
	d.HTTPHeader = httpheader.NewMiddleware(hh)
	d.TokenHandler = api.NewTokenHandler(tokenUsecase)
	d.MessageHandler = api.NewMessageHandler(messageUsecase)
}
