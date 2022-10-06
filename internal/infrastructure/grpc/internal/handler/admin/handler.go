package admin

import (
	admin_apiv1 "github.com/playground-live/moala-meet-and-greet-back/internal/infrastructure/grpc/pb/mmg/admin_api/v1"
	"github.com/playground-live/moala-meet-and-greet-back/internal/usecase"
)

type AdminHandler struct {
	admin_apiv1.UnimplementedAdminV1ServiceServer
	tenantInteractor usecase.AdminTenantInteractor
	userInteractor   usecase.AdminUserInteractor
}

func NewAdminHandler(
	tenantInteractor usecase.AdminTenantInteractor,
	userInteractor usecase.AdminUserInteractor,
) admin_apiv1.AdminV1ServiceServer {
	return &AdminHandler{
		tenantInteractor: tenantInteractor,
		userInteractor:   userInteractor,
	}
}
