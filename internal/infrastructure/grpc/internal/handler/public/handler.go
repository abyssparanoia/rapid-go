package public

import (
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database"
	public_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/public_api/v1"
)

type PublicHandler struct {
	databaseCli *database.Client
}

func NewPublicHandler(
	databaseCli *database.Client,
) public_apiv1.PublicV1ServiceServer {
	return &PublicHandler{
		databaseCli: databaseCli,
	}
}
