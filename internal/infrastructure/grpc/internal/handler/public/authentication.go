package public

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/interceptor/session_interceptor"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/internal/marshaller"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	public_apiv1 "github.com/abyssparanoia/rapid-go/schema/proto/pb/rapid/public_api/v1"
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
