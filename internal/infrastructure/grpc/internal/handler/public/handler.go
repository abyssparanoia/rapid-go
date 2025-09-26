package public

import (
	public_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/public_api/v1"
	database "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql"
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
