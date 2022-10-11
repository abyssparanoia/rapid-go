package public

import (
	"github.com/abyssparanoia/rapid-go/internal/usecase"
	public_apiv1 "github.com/abyssparanoia/rapid-go/schema/proto/pb/rapid/public_api/v1"
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
