package repository

import (
	"cloud.google.com/go/firestore"
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluefirestore"
	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/abyssparanoia/rapid-go/push-notification/domain/model"
	"github.com/abyssparanoia/rapid-go/push-notification/infrastructure/internal/entity"

	"context"
)

type token struct {
	firestoreClient *firestore.Client
}

func (r *token) GetByPlatformAndDeviceID(ctx context.Context,
	appID, userID, deviceID string,
	platform model.Platform) (*model.Token, error) {

	colRef := entity.NewTokenCollectionRef(r.firestoreClient, appID, userID)
	query := colRef.Where("device_id", "==", deviceID).Where("platform", "==", platform.String())

	tokenEntity := &entity.Token{}
	exist, err := gluefirestore.GetByQuery(ctx, query, tokenEntity)
	if err != nil {
		log.Errorm(ctx, "gluefirestore.GetByQuery", err)
		return nil, err
	}
	if !exist {
		return nil, nil
	}

	return tokenEntity.OutputModel(), nil
}
