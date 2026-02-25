//go:build gcp

// nolint:godot,gci
package dependency

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/firebase"
	firebase_repository "github.com/abyssparanoia/rapid-go/internal/infrastructure/firebase/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/gcs"
	gcs_repository "github.com/abyssparanoia/rapid-go/internal/infrastructure/gcs/repository"
	database "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql"
	database_cache "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql/cache"
	database_repository "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql/repository"
	database_transactable "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql/transactable"
	"github.com/abyssparanoia/rapid-go/internal/usecase"
)

type Dependency struct {
	DatabaseCli *database.Client

	// admin
	AdminTenantInteractor usecase.AdminTenantInteractor
	AdminStaffInteractor  usecase.AdminStaffInteractor
	AdminAssetInteractor  usecase.AdminAssetInteractor

	// staff
	StaffMeInteractor       usecase.StaffMeInteractor
	StaffMeTenantInteractor usecase.StaffMeTenantInteractor
	StaffStaffInteractor    usecase.StaffStaffInteractor
	StaffAssetInteractor    usecase.StaffAssetInteractor

	// Other
	AuthenticationInteractor      usecase.AuthenticationInteractor
	AdminAuthenticationInteractor usecase.AdminAuthenticationInteractor
	DebugInteractor               usecase.DebugInteractor

	// task
	TaskAdminInteractor usecase.TaskAdminInteractor
}

func (d *Dependency) Inject(
	ctx context.Context,
	e *environment.Environment,
) {
	d.DatabaseCli = database.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase, e.DBLogEnable)

	// GCP Cloud Storage
	gcsCli := gcs.NewClient(ctx, e.GCSEmulatorHost)
	gcsPrivateBucketHandle := gcs.NewBucketHandle(gcsCli, e.GCPPrivateBucketName)
	gcsPublicBucketHandle := gcs.NewBucketHandle(gcsCli, e.GCPPublicBucketName)

	// Firebase Auth
	firebaseCli := firebase.NewClient(e.GCPProjectID, e.FirebaseAuthEmulatorHost)

	transactable := database_transactable.NewTransactable()

	staffAuthenticationRepository := firebase_repository.NewStaffAuthentication(
		firebaseCli,
		e.FirebaseClientAPIKey,
		e.FirebaseAuthEmulatorHost,
	)
	adminAuthenticationRepository := firebase_repository.NewAdminAuthentication(
		firebaseCli,
		e.FirebaseClientAPIKey,
		e.FirebaseAuthEmulatorHost,
	)
	tenantRepository := database_repository.NewTenant()
	staffRepository := database_repository.NewStaff()
	adminRepository := database_repository.NewAdmin()

	// GCS asset repository
	assetRepository := gcs_repository.NewAsset(
		gcsPrivateBucketHandle,
		gcsPublicBucketHandle,
		e.GCPPublicAssetBaseURL,
		e.GCSEmulatorHost,
		e.GCPPrivateBucketName,
		e.GCPPublicBucketName,
	)

	assetPathCache := database_cache.NewAssetPath()

	assetService := service.NewAsset(
		assetRepository,
		assetPathCache,
	)

	staffService := service.NewStaff(
		staffRepository,
		staffAuthenticationRepository,
	)

	d.AdminTenantInteractor = usecase.NewAdminTenantInteractor(
		transactable,
		tenantRepository,
		assetService,
	)
	d.AdminStaffInteractor = usecase.NewAdminStaffInteractor(
		transactable,
		tenantRepository,
		staffRepository,
		staffService,
		assetService,
	)
	d.AdminAssetInteractor = usecase.NewAdminAssetInteractor(
		assetService,
	)

	d.StaffMeInteractor = usecase.NewStaffMeInteractor(
		transactable,
		tenantRepository,
		staffRepository,
		staffService,
		assetService,
	)
	d.StaffMeTenantInteractor = usecase.NewStaffMeTenantInteractor(
		transactable,
		tenantRepository,
		assetService,
	)
	d.StaffStaffInteractor = usecase.NewStaffStaffInteractor(
		transactable,
		tenantRepository,
		staffRepository,
		staffService,
		assetService,
	)
	d.StaffAssetInteractor = usecase.NewStaffAssetInteractor(
		assetService,
	)

	d.AuthenticationInteractor = usecase.NewAuthenticationInteractor(
		staffAuthenticationRepository,
	)

	d.AdminAuthenticationInteractor = usecase.NewAdminAuthenticationInteractor(
		adminAuthenticationRepository,
	)

	d.DebugInteractor = usecase.NewDebugInteractor(
		adminAuthenticationRepository,
		staffAuthenticationRepository,
	)

	d.TaskAdminInteractor = usecase.NewTaskAdminInteractor(
		transactable,
		adminRepository,
		adminAuthenticationRepository,
	)
}
