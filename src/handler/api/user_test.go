package api

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abyssparanoia/rapid-go/src/domain/model"
	mock_service "github.com/abyssparanoia/rapid-go/src/service/mock"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockExpectedGet struct {
	userID int64
	result *model.User
	err    error
}

func TestUserHandler_Get(t *testing.T) {
	tests := []struct {
		name   string
		mock   mockExpectedGet
		param  map[string]string
		want   string
		status int
	}{
		{
			name: "success",
			mock: mockExpectedGet{
				userID: int64(1),
				result: &model.User{
					ID:   1,
					Name: "abyssparanoia",
					Sex:  "man",
				},
				err: nil,
			},
			param: map[string]string{"userID": "1"},
			want: `{
				"user": {
					"id": 1,
					"name": "abyssparanoia",
					"sex" : "man"
				}
			}`,
			status: http.StatusOK,
		},
		{
			name:   "invalid tyeps userid",
			mock:   mockExpectedGet{},
			param:  map[string]string{"userID": "a"},
			want:   ``,
			status: http.StatusInternalServerError,
		},
		{
			name: "service error",
			mock: mockExpectedGet{
				userID: int64(1),
				result: &model.User{},
				err:    errors.New("error"),
			},
			param:  map[string]string{"userID": "1"},
			want:   ``,
			status: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mc := gomock.NewController(t)
			mockUser := mock_service.NewMockUser(mc)
			mockUser.EXPECT().Get(gomock.Any(), tt.mock.userID).Return(tt.mock.result, tt.mock.err)
			handler := NewUserHandler(mockUser)

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.Get(w, r)
			})

			r, _ := http.NewRequest("GET", "/", nil)
			rctx := chi.NewRouteContext()
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			for k, v := range tt.param {
				rctx.URLParams.Add(k, v)
			}

			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)

			resp := w.Result()
			if tt.status != http.StatusOK {
				return
			}
			got, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if !assert.JSONEq(t, tt.want, string(got)) {
				t.Errorf("UserHandler.Get() = %v, want %v", string(got), tt.want)
			}
		})

	}
}
