package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/abyssparanoia/rapid-go/internal/pkg/gluefirestore"
	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/abyssparanoia/rapid-go/push-notification/domain/model"
	"github.com/abyssparanoia/rapid-go/push-notification/domain/repository"
	"github.com/abyssparanoia/rapid-go/push-notification/infrastructure/internal/entity"
	"google.golang.org/api/iterator"
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

func (r *token) List(ctx context.Context) ([]*model.Token, error) {

	tokenEntityList := []*entity.Token{}
	var cursor *firestore.DocumentSnapshot
	for {
		_tokenEntityList := []*entity.Token{}
		var dsnp *firestore.DocumentSnapshot
		query := entity.NewTokenGroupCollectionRef(r.firestoreClient)
		if cursor != nil {
			query.StartAfter(cursor)
		}
		it := query.Documents(ctx)
		for {
			dsnp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Errorm(ctx, "it.Next", err)
				return nil, err
			}
			tokenEntity := &entity.Token{}
			err = dsnp.DataTo(tokenEntity)
			if err != nil {
				log.Errorm(ctx, "doc.DataTo", err)
				return nil, err
			}
			_tokenEntityList = append(_tokenEntityList, tokenEntity)
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

	colRef := entity.NewTokenCollectionRef(r.firestoreClient, appID, userID)

	tokenEntityList := []*entity.Token{}
	err := gluefirestore.ListByQuery(ctx, colRef.Query, tokenEntityList)
	if err != nil {
		log.Errorm(ctx, "gluefirestore.ListByQuery", err)
		return nil, err
	}
	return entity.NewTokenMultiOutputModels(tokenEntityList), nil
}

// NewToken ... new token repository
func NewToken(firestoreClient *firestore.Client) repository.Token {
	return &token{firestoreClient}
}
