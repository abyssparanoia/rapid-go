package api

import (
	"net/http"
	"strconv"

	"github.com/abyssparanoia/rapid-go/src/domain/model"
	"github.com/abyssparanoia/rapid-go/src/handler"
	"github.com/abyssparanoia/rapid-go/src/service"
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
		handler.HandleError(ctx, w, "validation error: ", err)
		return
	}
	param.UserID = userID

	v := validator.New()
	if err = v.Struct(param); err != nil {
		handler.HandleError(ctx, w, "validation error: ", err)
		return
	}

	user, err := h.Svc.Get(ctx, param.UserID)
	if err != nil {
		handler.HandleError(ctx, w, "h.Svc.Get: ", err)
	}

	handler.RenderJSON(w, http.StatusOK, userHandlerGetResponse{User: user})
}

// NewUserHandler ... ユーザーハンドラーを取得する
func NewUserHandler(Svc service.User) *UserHandler {
	return &UserHandler{
		Svc: Svc,
	}
}
