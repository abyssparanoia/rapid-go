package public

import (
	"context"

	"github.com/playground-live/moala-meet-and-greet-back/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
	"github.com/playground-live/moala-meet-and-greet-back/internal/infrastructure/grpc/internal/marshaller"
	public_apiv1 "github.com/playground-live/moala-meet-and-greet-back/internal/infrastructure/grpc/pb/mmg/public_api/v1"
	"github.com/playground-live/moala-meet-and-greet-back/internal/usecase/input"
)

func (h *PublicHandler) SignIn(ctx context.Context, req *public_apiv1.SignInRequest) (*public_apiv1.SignInResponse, error) {
	sctx, err := session_interceptor.RequireSessionContext(ctx)
	if err != nil {
		return nil, err
	}
	got, err := h.authenticationInteractor.SignIn(
		ctx,
		input.NewPublicSignIn(
			sctx.AuthUID,
		),
	)
	if err != nil {
		return nil, err
	}

	return &public_apiv1.SignInResponse{
		User: marshaller.UserToPB(got),
	}, nil
}
