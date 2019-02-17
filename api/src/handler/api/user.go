package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/abyssparanoia/rapid-go/api/src/handler"
	"github.com/abyssparanoia/rapid-go/api/src/lib/log"
	"github.com/abyssparanoia/rapid-go/api/src/model"
	"github.com/abyssparanoia/rapid-go/api/src/service"
	validator "gopkg.in/go-playground/validator.v9"
)

// UserHandler ... ユーザーハンドラー
type UserHandler struct {
	Svc service.User
}

type userHandlerGetRequestParam struct {
	UserID int64 `validate:"required"`
}

type userHandlerGetResponse struct {
	User *model.User `json:"user"`
}

// Get ... ユーザー情報を取得する
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var param userHandlerGetRequestParam
	userID, err := strconv.ParseInt(handler.GetURLParam(r, "userID"), 10, 64)
	if err != nil {
		h.handleError(ctx, w, http.StatusBadRequest, "validation error: "+err.Error())
		return
	}
	param.UserID = userID

	v := validator.New()
	if err = v.Struct(param); err != nil {
		h.handleError(ctx, w, http.StatusBadRequest, "validation error: "+err.Error())
		return
	}

	user, err := h.Svc.Get(ctx, param.UserID)
	if err != nil {
		h.handleError(ctx, w, http.StatusInternalServerError, "h.Svc.Get: "+err.Error())
	}

	handler.RenderJSON(w, http.StatusOK, userHandlerGetResponse{User: user})
}

func (h *UserHandler) handleError(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	log.Errorf(ctx, msg)
	handler.RenderError(w, status, msg)
}

// NewUserHandler ... ユーザーハンドラーを取得する
func NewUserHandler(Svc service.User) *UserHandler {
	return &UserHandler{
		Svc: Svc,
	}
}
