package api

import (
	"net/http"

	"github.com/abyssparanoia/rapid-go/internal/default/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/default/usecase"
	"github.com/abyssparanoia/rapid-go/internal/pkg/parameter"
	"github.com/abyssparanoia/rapid-go/internal/pkg/renderer"
	validator "github.com/go-playground/validator"
)

// UserHandler ... user handler
type UserHandler struct {
	userUsecase usecase.User
}

type userHandlerGetRequestParam struct {
	UserID string `validate:"required"`
}

type userHandlerGetResponse struct {
	User *model.User `json:"user"`
}

// Get ... get user
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var param userHandlerGetRequestParam
	param.UserID = parameter.GetURL(r, "userID")

	v := validator.New()
	if err := v.Struct(param); err != nil {
		renderer.HandleError(ctx, w, err)
		return
	}

	user, err := h.userUsecase.Get(ctx, param.UserID)
	if err != nil {
		renderer.HandleError(ctx, w, err)
		return
	}

	renderer.JSON(ctx, w, http.StatusOK, userHandlerGetResponse{User: user})
}

// NewUserHandler ... get user handler
func NewUserHandler(userUsecase usecase.User) *UserHandler {
	return &UserHandler{userUsecase}
}
