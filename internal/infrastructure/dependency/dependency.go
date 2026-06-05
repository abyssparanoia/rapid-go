//go:build !gcp

// nolint:godot,gci
package dependency

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/aws"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito"
	cognito_repository "github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	database "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql"
	database_cache "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql/cache"
	database_repository "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql/repository"
	database_transactable "github.com/abyssparanoia/rapid-go/internal/infrastructure/mysql/transactable"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/s3"
	s3_repository "github.com/abyssparanoia/rapid-go/internal/infrastructure/s3/repository"
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

	// AWS
	awsSession := aws.NewConfig(ctx, e.AWSRegion)
	s3Client := s3.NewClient(awsSession, e.AWSEmulatorHost)
	cognitoCli := cognito.NewClient(awsSession, e.AWSCognitoEmulatorHost)

	transactable := database_transactable.NewTransactable()

	staffAuthenticationRepository := cognito_repository.NewStaffAuthentication(
		ctx,
		cognitoCli,
		e.AWSCognitoStaffUserPoolID,
		e.AWSCognitoStaffClientID,
		e.AWSCognitoEmulatorHost,
		e.AWSRegion,
	)
	adminAuthenticationRepository := cognito_repository.NewAdminAuthentication(
		ctx,
		cognitoCli,
		e.AWSCognitoAdminUserPoolID,
		e.AWSCognitoAdminClientID,
		e.AWSCognitoEmulatorHost,
		e.AWSRegion,
	)
	tenantRepository := database_repository.NewTenant()
	staffRepository := database_repository.NewStaff()
	adminRepository := database_repository.NewAdmin()

	// S3 asset repository
	assetRepository := s3_repository.NewAsset(
		s3Client,
		e.AWSPrivateBucketName,
		e.AWSPublicBucketName,
		e.AWSPublicAssetBaseURL,
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
		staffRepository,
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
