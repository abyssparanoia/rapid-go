package admin

import (
	"github.com/abyssparanoia/rapid-go/internal/usecase"
	admin_apiv1 "github.com/abyssparanoia/rapid-go/schema/proto/pb/rapid/admin_api/v1"
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
