package firebaseauth

import (
	"context"
	"net/http"
)

type dummyService struct {
}

// SetCustomClaims ... カスタムClaimsを設定
func (s *dummyService) SetCustomClaims(ctx context.Context, userID string, claims Claims) error {
	return nil
}

// Authentication ... 認証を行う
func (s *dummyService) Authentication(ctx context.Context, r *http.Request) (string, Claims, error) {
	userID := "DUMMY_USER_ID"
	claims := Claims{
		// EDIT: 任意でClaimsにダミーデータを入れる
	}
	return userID, claims, nil
}

// NewDummyService ... DummyServiceを作成する
func NewDummyService() Service {
	return &dummyService{}
}
