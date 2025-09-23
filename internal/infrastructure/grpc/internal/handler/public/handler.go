package public

import (
	public_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/public_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql"
)

type PublicHandler struct {
	databaseCli *mysql.Client
}

func NewPublicHandler(
	databaseCli *mysql.Client,
) public_apiv1.PublicV1ServiceServer {
	return &PublicHandler{
		databaseCli: databaseCli,
	}
}
