package staff

import (
	staff_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/staff_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase"
)

type StaffHandler struct {
	tenantInteractor usecase.StaffTenantInteractor
	staffInteractor  usecase.StaffStaffInteractor
	assetInteractor  usecase.StaffAssetInteractor
}

func NewStaffHandler(
	tenantInteractor usecase.StaffTenantInteractor,
	staffInteractor usecase.StaffStaffInteractor,
	assetInteractor usecase.StaffAssetInteractor,
) staff_apiv1.StaffV1ServiceServer {
	return &StaffHandler{
		tenantInteractor: tenantInteractor,
		staffInteractor:  staffInteractor,
		assetInteractor:  assetInteractor,
	}
}
