package service

import (
	"context"
	"errors"
	"testing"

	"github.com/abyssparanoia/rapid-go/src/domain/model"
	mock_repository "github.com/abyssparanoia/rapid-go/src/domain/repository/mock"
	"github.com/abyssparanoia/rapid-go/src/infrastructure/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockExpectedGet struct {
	userID int64
	result *entity.User
	err    error
}

func Test_user_Get(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID int64
	}

	tests := []struct {
		name    string
		mock    mockExpectedGet
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "success sample",
			mock: mockExpectedGet{
				userID: 1,
				result: &entity.User{
					ID:   1,
					Name: "abyssparanoia",
					Sex:  "man",
				},
				err: nil,
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want: &model.User{
				ID:   1,
				Name: "abyssparanoia",
			},
		},
		{
			name: "error sample",
			mock: mockExpectedGet{
				userID: 0,
				result: &entity.User{},
				err:    errors.New("test"),
			},
			args: args{
				ctx:    context.Background(),
				userID: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			uRepo := mock_repository.NewMockUser(mc)
			s := NewUser(uRepo)
			got, err := s.Get(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("user.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
