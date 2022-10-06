package public

import (
	public_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/mmg/public_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase"
)

type PublicHandler struct {
	public_apiv1.UnimplementedPublicV1ServiceServer
	authenticationInteractor usecase.PublicAuthenticationInteractor
	tenantInteractor         usecase.PublicTenantInteractor
}

func NewPublicHandler(
	authenticationInteractor usecase.PublicAuthenticationInteractor,
	tenantInteractor usecase.PublicTenantInteractor,
) public_apiv1.PublicV1ServiceServer {
	return &PublicHandler{
		authenticationInteractor: authenticationInteractor,
		tenantInteractor:         tenantInteractor,
	}
}
