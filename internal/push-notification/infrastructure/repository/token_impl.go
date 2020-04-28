package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluefirestore"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/infrastructure/internal/entity"
)

type token struct {
	firestoreClient *firestore.Client
}

func (r *token) GetByPlatformAndDeviceIDAndUserID(ctx context.Context,
	appID, userID, deviceID string,
	platform model.Platform) (*model.Token, error) {

	colRef := entity.NewTokenCollectionRef(r.firestoreClient)
	query := colRef.
		Where("app_id", "==", appID).
		Where("user_id", "==", userID).
		Where("device_id", "==", deviceID).
		Where("platform", "==", platform.String())

	tokenEntity := &entity.Token{}
	exist, err := gluefirestore.GetByQuery(ctx, query, tokenEntity)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}

	return tokenEntity.OutputModel(), nil
}

func (r *token) List(ctx context.Context) ([]*model.Token, error) {

	tokenEntityList := []*entity.Token{}
	var cursor *firestore.DocumentSnapshot

	colRef := entity.NewTokenCollectionRef(r.firestoreClient)

	for {
		_tokenEntityList := []*entity.Token{}
		var dsnp *firestore.DocumentSnapshot
		query := colRef.Query

		if cursor != nil {
			query.StartAfter(cursor)
		}

		err := gluefirestore.ListByQuery(ctx, query, _tokenEntityList)
		if err != nil {
			return nil, err
		}

		var nCursor *firestore.DocumentSnapshot
		if len(_tokenEntityList) == 300 {
			nCursor = dsnp
		}
		for _, tokenEntity := range _tokenEntityList {
			tokenEntityList = append(tokenEntityList, tokenEntity)
		}
		if nCursor == nil {
			break
		}
		cursor = nCursor
	}

	return entity.NewTokenMultiOutputModels(tokenEntityList), nil
}

func (r *token) ListByUserID(ctx context.Context,
	appID, userID string) ([]*model.Token, error) {

	colRef := entity.NewTokenCollectionRef(r.firestoreClient)
	query := colRef.
		Where("app_id", "==", appID).
		Where("user_id", "==", userID)

	tokenEntityList := []*entity.Token{}
	err := gluefirestore.ListByQuery(ctx, query, tokenEntityList)
	if err != nil {
		return nil, err
	}
	return entity.NewTokenMultiOutputModels(tokenEntityList), nil
}

func (r *token) Create(ctx context.Context,
	token *model.Token) (*model.Token, error) {

	tokenEntity := entity.NewTokenFromModel(token)

	colRef := entity.NewTokenCollectionRef(r.firestoreClient)

	err := gluefirestore.Create(ctx, colRef, tokenEntity)
	if err != nil {
		return nil, err
	}

	return tokenEntity.OutputModel(), nil
}

func (r *token) Update(ctx context.Context,
	token *model.Token) error {

	tokenEntity := entity.NewTokenFromModel(token)

	colRef := entity.NewTokenCollectionRef(r.firestoreClient)
	docRef := colRef.Doc(token.ID)
	err := gluefirestore.Set(ctx, docRef, tokenEntity)
	if err != nil {
		return err
	}

	return nil
}

func (r *token) Delete(ctx context.Context,
	tokenID string) error {

	colRef := entity.NewTokenCollectionRef(r.firestoreClient)
	docRef := colRef.Doc(tokenID)

	err := gluefirestore.Delete(ctx, docRef)
	if err != nil {
		return err
	}

	return nil
}

// NewToken ... new token repository
func NewToken(firestoreClient *firestore.Client) repository.Token {
	return &token{firestoreClient}
}
