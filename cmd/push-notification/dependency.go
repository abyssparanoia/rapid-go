package main

import (
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluefcm"
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluefirestore"
	"github.com/abyssparanoia/rapid-go/internal/pkg/httpheader"
	"github.com/abyssparanoia/rapid-go/internal/pkg/httpmiddleware"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/handler/api"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/infrastructure/repository"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/usecase"
	"go.uber.org/zap"
)

// Dependency ... dependency
type Dependency struct {
	httpMiddleware  *httpmiddleware.HTTPMiddleware
	DummyHTTPHeader *httpheader.Middleware
	HTTPHeader      *httpheader.Middleware
	TokenHandler    *api.TokenHandler
	MessageHandler  *api.MessageHandler
}

// Inject ... inject dependency
func (d *Dependency) Inject(e *environment, logger *zap.Logger) {

	fcmClient := gluefcm.NewClient(e.ProjectID)
	firestoreClient := gluefirestore.NewClient(e.ProjectID)

	tokenRepository := repository.NewToken(firestoreClient)
	fcmRepository := repository.NewFcm(fcmClient, e.FcmServerKey)

	tokenService := service.NewToken(tokenRepository)

	tokenUsecase := usecase.NewToken(fcmRepository, tokenRepository, tokenService)
	messageUsecase := usecase.NewMessage(fcmRepository, tokenRepository)

	dhh := httpheader.NewDummy()
	hh := httpheader.New()

	d.httpMiddleware = httpmiddleware.New(logger)

	d.DummyHTTPHeader = httpheader.NewMiddleware(dhh)
	d.HTTPHeader = httpheader.NewMiddleware(hh)
	d.TokenHandler = api.NewTokenHandler(tokenUsecase)
	d.MessageHandler = api.NewMessageHandler(messageUsecase)
}
