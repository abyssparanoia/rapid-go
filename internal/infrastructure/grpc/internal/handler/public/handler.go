package public

import (
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database"
	public_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/public_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase"
)

type PublicHandler struct {
	databaseCli      *database.Client
	tenantInteractor usecase.PublicTenantInteractor
}

func NewPublicHandler(
	databaseCli *database.Client,
	tenantInteractor usecase.PublicTenantInteractor,
) public_apiv1.PublicV1ServiceServer {
	return &PublicHandler{
		databaseCli:      databaseCli,
		tenantInteractor: tenantInteractor,
	}
}
