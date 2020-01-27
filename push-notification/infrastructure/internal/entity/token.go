package entity

import (
	"cloud.google.com/go/firestore"

	"github.com/abyssparanoia/rapid-go/push-notification/domain/model"
)

// Token ... token entity
type Token struct {
	ID        string                 `firestore:"-" gluefirestore:"id"`
	Ref       *firestore.DocumentRef `firestore:"-" gluefirestore:"ref"`
	AppID     string                 `firestore:"app_id"`
	UserID    string                 `firestore:"user_id"`
	Platform  string                 `firestore:"platform"`
	DeviceID  string                 `firestore:"device_id"`
	Value     string                 `firestore:"value"`
	CreatedAt int64                  `firestore:"created_at"`
}

// BuildFromModel ... build from model
func (e *Token) BuildFromModel(m *model.Token) {
	e.ID = m.ID
	e.AppID = m.AppID
	e.UserID = m.UserID
	e.Platform = m.Platform.String()
	e.DeviceID = m.DeviceID
	e.Value = m.Value
	e.CreatedAt = m.CreatedAt
}

// OutputModel ... output model
func (e *Token) OutputModel() *model.Token {
	return &model.Token{
		ID:        e.ID,
		AppID:     e.AppID,
		UserID:    e.UserID,
		Platform:  model.MustPlatform(e.Platform),
		DeviceID:  e.DeviceID,
		Value:     e.Value,
		CreatedAt: e.CreatedAt,
	}
}

// NewTokenMultiOutputModels ... multi output models
func NewTokenMultiOutputModels(dsts []*Token) (tokens []*model.Token) {
	for _, dst := range dsts {
		tokens = append(tokens, dst.OutputModel())
	}
	return tokens
}

// NewTokenCollectionRef ... new token collection ref
func NewTokenCollectionRef(fCli *firestore.Client) *firestore.CollectionRef {
	return fCli.Collection("push-notification-tokens")
}
