package public

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/marshaller"
	public_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/public_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
)

func (h *PublicHandler) PublicSignIn(ctx context.Context, req *public_apiv1.PublicSignInRequest) (*public_apiv1.PublicSignInResponse, error) {
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

	return &public_apiv1.PublicSignInResponse{
		User: marshaller.UserToPB(got),
	}, nil
}
