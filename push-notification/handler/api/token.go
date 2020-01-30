package api

import (
	"net/http"

	"github.com/abyssparanoia/rapid-go/push-notification/usecase/input"

	"github.com/abyssparanoia/rapid-go/internal/pkg/errcode"
	"github.com/abyssparanoia/rapid-go/internal/pkg/parameter"
	"github.com/abyssparanoia/rapid-go/internal/pkg/renderer"
	"github.com/abyssparanoia/rapid-go/push-notification/usecase"
)

// TokenHandler ... token handler struct
type TokenHandler struct {
	tokenUsecase usecase.Token
}

// Set ... token set handler
func (h *TokenHandler) Set(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var param struct {
		AppID    string `json:"app_id" validate:"required"`
		UserID   string `json:"user_id" validate:"required"`
		Platform string `json:"platform" validate:"required"`
		DeviceID string `json:"device_id" validate:"required"`
		Token    string `json:"token" validate:"required"`
	}

	err := parameter.GetJSON(r, &param)
	if err != nil {
		err = errcode.Set(err, http.StatusBadRequest)
		renderer.HandleError(ctx, w, "paramater.GetJSON", err)
		return
	}

	dto, err := input.NewTokenSet(ctx, param.AppID, param.UserID, param.Platform, param.DeviceID, param.Token)
	if err != nil {
		renderer.HandleError(ctx, w, "input.NewTokenSet", err)
		return
	}

	err = h.tokenUsecase.Set(ctx, dto)
	if err != nil {
		renderer.HandleError(ctx, w, "h.tokenUsecase.Set", err)
		return
	}

	renderer.Success(ctx, w)
}

// Delete ... delete token handler
func (h *TokenHandler) Delete(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var param struct {
		AppID    string `json:"app_id" validate:"required"`
		UserID   string `json:"user_id" validate:"required"`
		Platform string `json:"platform" validate:"required"`
		DeviceID string `json:"device_id" validate:"required"`
	}

	err := parameter.GetJSON(r, &param)
	if err != nil {
		err = errcode.Set(err, http.StatusBadRequest)
		renderer.HandleError(ctx, w, "paramater.GetJSON", err)
		return
	}

	dto, err := input.NewTokenDelete(ctx, param.AppID, param.UserID, param.Platform, param.DeviceID)
	if err != nil {
		renderer.HandleError(ctx, w, "input.NewTokenDelete", err)
		return
	}

	err = h.tokenUsecase.Delete(ctx, dto)
	if err != nil {
		renderer.HandleError(ctx, w, "h.tokenUsecase.Delete", err)
		return
	}

	renderer.Success(ctx, w)
}

// NewTokenHandler ... new token handler
func NewTokenHandler(tokenUsecase usecase.Token) *TokenHandler {
	return &TokenHandler{tokenUsecase}
}
