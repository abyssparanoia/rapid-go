package staff

import (
	staff_apiv1 "github.com/abyssparanoia/rapid-go/internal/infrastructure/grpc/pb/rapid/staff_api/v1"
	"github.com/abyssparanoia/rapid-go/internal/usecase"
)

type StaffHandler struct {
	meInteractor       usecase.StaffMeInteractor
	meTenantInteractor usecase.StaffMeTenantInteractor
	staffInteractor    usecase.StaffStaffInteractor
	assetInteractor    usecase.StaffAssetInteractor
}

func NewStaffHandler(
	meInteractor usecase.StaffMeInteractor,
	meTenantInteractor usecase.StaffMeTenantInteractor,
	staffInteractor usecase.StaffStaffInteractor,
	assetInteractor usecase.StaffAssetInteractor,
) staff_apiv1.StaffV1ServiceServer {
	return &StaffHandler{
		meInteractor:       meInteractor,
		meTenantInteractor: meTenantInteractor,
		staffInteractor:    staffInteractor,
		assetInteractor:    assetInteractor,
	}
}
