package admin

import (
	admin_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/admin_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase"
)

type AdminHandler struct {
	tenantInteractor usecase.AdminTenantInteractor
	staffInteractor  usecase.AdminStaffInteractor
	assetInteractor  usecase.AdminAssetInteractor
}

func NewAdminHandler(
	tenantInteractor usecase.AdminTenantInteractor,
	staffInteractor usecase.AdminStaffInteractor,
	assetInteractor usecase.AdminAssetInteractor,
) admin_apiv1.AdminV1ServiceServer {
	return &AdminHandler{
		tenantInteractor: tenantInteractor,
		staffInteractor:  staffInteractor,
		assetInteractor:  assetInteractor,
	}
}
