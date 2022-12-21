package admin

import (
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase"
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
