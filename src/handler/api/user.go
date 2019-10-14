package api

import (
	"net/http"

	"github.com/abyssparanoia/rapid-go/src/domain/model"
	"github.com/abyssparanoia/rapid-go/src/lib/parameter"
	"github.com/abyssparanoia/rapid-go/src/lib/renderer"
	"github.com/abyssparanoia/rapid-go/src/service"
	validator "gopkg.in/go-playground/validator.v9"
)

// UserHandler ... user handler
type UserHandler struct {
	Svc service.User
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
		renderer.HandleError(ctx, w, "validation error: ", err)
		return
	}

	user, err := h.Svc.Get(ctx, param.UserID)
	if err != nil {
		renderer.HandleError(ctx, w, "h.Svc.Get: ", err)
	}

	renderer.JSON(ctx, w, http.StatusOK, userHandlerGetResponse{User: user})
}

// NewUserHandler ... get user handler
func NewUserHandler(Svc service.User) *UserHandler {
	return &UserHandler{
		Svc: Svc,
	}
}
