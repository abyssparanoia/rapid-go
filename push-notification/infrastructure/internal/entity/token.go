package entity

import (
	"cloud.google.com/go/firestore"

	"github.com/abyssparanoia/rapid-go/push-notification/domain/model"
)

// Token ... token entity
type Token struct {
	ID        string                 `firestore:"-" gluefirestore:"id"`
	Ref       *firestore.DocumentRef `firestore:"-" gluefirestore:"ref"`
	Platform  string                 `firestore:"platform"`
	DeviceID  string                 `firestore:"device_id"`
	Token     string                 `firestore:"token"`
	CreatedAt int64                  `firestore:"created_at"`
}

// BuildFromModel ... build from model
func (e *Token) BuildFromModel(m *model.Token) {
	e.ID = m.ID
	e.Platform = m.Platform.String()
	e.DeviceID = m.DeviceID
	e.Token = m.Token
	e.CreatedAt = m.CreatedAt
}

// OutputModel ... output model
func (e *Token) OutputModel() *model.Token {
	return &model.Token{
		ID:        e.ID,
		Platform:  model.MustPlatform(e.Platform),
		DeviceID:  e.DeviceID,
		Token:     e.Token,
		CreatedAt: e.CreatedAt,
	}
}

// NewTokenCollectionRef ... new token collection ref
func NewTokenCollectionRef(fCli *firestore.Client, appID, userID string) *firestore.CollectionRef {
	return fCli.Collection("push-notification-apps").Doc(appID).Collection("users").Doc(userID).Collection("tokens")
}
