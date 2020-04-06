package api

import (
	"net/http"

	"github.com/abyssparanoia/rapid-go/internal/pkg/errcode"
	"github.com/abyssparanoia/rapid-go/internal/pkg/parameter"
	"github.com/abyssparanoia/rapid-go/internal/pkg/renderer"
	"github.com/abyssparanoia/rapid-go/internal/push-notification/usecase/input"

	"github.com/abyssparanoia/rapid-go/internal/push-notification/usecase"
	validator "github.com/go-playground/validator"
)

// MessageHandler ... message handler
type MessageHandler struct {
	messageUsecase usecase.Message
}

// SendToUser ... send to user handler
func (h *MessageHandler) SendToUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var param struct {
		AppID          string                `json:"app_id" validate:"required"`
		UserID         string                `json:"user_id" validate:"required"`
		MessageRequest *input.MessageRequest `json:"message" validate:"required"`
	}

	err := parameter.GetJSON(r, &param)
	if err != nil {
		err = errcode.Set(err, http.StatusBadRequest)
		renderer.HandleError(ctx, w, "paramater.GetJSON", err)
		return
	}

	v := validator.New()
	if err := v.Struct(param); err != nil {
		renderer.HandleError(ctx, w, "validation error: ", err)
		return
	}

	dto := input.NewMessageSendToUser(param.AppID, param.UserID, param.MessageRequest)
	err = h.messageUsecase.SendToUser(ctx, dto)
	if err != nil {
		renderer.HandleError(ctx, w, "h.messageUsecase.SendToUser", err)
		return
	}

	renderer.Success(ctx, w)
}

// SendToMultiUser ... send to multi user handler
func (h *MessageHandler) SendToMultiUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var param struct {
		AppID          string                `json:"app_id" validate:"required"`
		UserIDList     []string              `json:"user_id_list" validate:"required"`
		MessageRequest *input.MessageRequest `json:"message" validate:"required"`
	}

	err := parameter.GetJSON(r, &param)
	if err != nil {
		err = errcode.Set(err, http.StatusBadRequest)
		renderer.HandleError(ctx, w, "paramater.GetJSON", err)
		return
	}

	v := validator.New()
	if err := v.Struct(param); err != nil {
		renderer.HandleError(ctx, w, "validation error: ", err)
		return
	}

	dto := input.NewMessageSendToMultiUser(param.AppID, param.UserIDList, param.MessageRequest)
	err = h.messageUsecase.SendToMultiUser(ctx, dto)
	if err != nil {
		renderer.HandleError(ctx, w, "h.messageUsecase.SendToMultiUser", err)
		return
	}

	renderer.Success(ctx, w)
}

// SendToAllUser ... send to all user handler
func (h *MessageHandler) SendToAllUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var param struct {
		AppID          string                `json:"app_id" validate:"required"`
		MessageRequest *input.MessageRequest `json:"message" validate:"required"`
	}

	err := parameter.GetJSON(r, &param)
	if err != nil {
		err = errcode.Set(err, http.StatusBadRequest)
		renderer.HandleError(ctx, w, "paramater.GetJSON", err)
		return
	}

	v := validator.New()
	if err := v.Struct(param); err != nil {
		renderer.HandleError(ctx, w, "validation error: ", err)
		return
	}

	dto := input.NewMessageSendToAllUser(param.AppID, param.MessageRequest)
	err = h.messageUsecase.SendToAllUser(ctx, dto)
	if err != nil {
		renderer.HandleError(ctx, w, "h.messageUsecase.SendToAllUser", err)
		return
	}

	renderer.Success(ctx, w)
}

// NewMessageHandler ... new message handler
func NewMessageHandler(messageUsecase usecase.Message) *MessageHandler {
	return &MessageHandler{messageUsecase}
}
