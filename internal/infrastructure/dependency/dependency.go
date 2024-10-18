// nolint:godot,gci
package dependency

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/aws"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito"
	cognito_repository "github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database"
	database_cache "github.com/abyssparanoia/rapid-go/internal/infrastructure/database/cache"

	// redis_cache "github.com/abyssparanoia/rapid-go/internal/infrastructure/redis/cache"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/transactable"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/firebase"
	firebase_repository "github.com/abyssparanoia/rapid-go/internal/infrastructure/firebase/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/gcs"
	gcs_repository "github.com/abyssparanoia/rapid-go/internal/infrastructure/gcs/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/s3"
	"github.com/abyssparanoia/rapid-go/internal/usecase"
)

type Dependency struct {
	DatabaseCli *database.Client

	// admin
	AdminTenantInteractor usecase.AdminTenantInteractor
	AdminStaffInteractor  usecase.AdminStaffInteractor
	AdminAssetInteractor  usecase.AdminAssetInteractor

	// Other
	StaffInteractor          usecase.StaffInteractor
	AuthenticationInteractor usecase.AuthenticationInteractor
	DebugInteractor          usecase.DebugInteractor
}

func (d *Dependency) Inject(
	ctx context.Context,
	e *environment.Environment,
) {
	d.DatabaseCli = database.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase, e.DBLogEnable)
	// redisCli := redis.NewClient(e.RedisHost, e.RedisPort, e.RedisUsername, e.RedisPassword, e.RedisTLSEnable)

	firebaseCli := firebase.NewClient(e.GCPProjectID)

	gcsCli := gcs.NewClient(ctx)
	gcsBucketHandle := gcs.NewBucketHandle(gcsCli, e.GCPBucketName)

	awsSession := aws.NewConfig(ctx, e.AWSRegion)
	_ = s3.NewClient(awsSession, e.AWSEmulatorHost)

	cognitoCli := cognito.NewClient(awsSession, e.AWSCognitoEmulatorHost)

	transactable := transactable.NewTransactable()
	_ = firebase_repository.NewStaffAuthentication(
		firebaseCli,
		e.FirebaseClientAPIKey,
	)
	staffAuthenticationRepository := cognito_repository.NewStaffAuthentication(
		ctx,
		cognitoCli,
		e.AWSCognitoUserPoolID,
		e.AWSCognitoClientID,
		e.AWSCognitoEmulatorHost,
	)
	tenantRepository := repository.NewTenant()
	staffRepository := repository.NewStaff()
	assetRepository := gcs_repository.NewAsset(gcsBucketHandle)
	// assetRepository := s3_repository.NewAsset(s3Client, e.AWSBucketName)

	// assetPathCache := redis_cache.NewAssetPath(redisCli)
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
	)
	d.AdminStaffInteractor = usecase.NewAdminStaffInteractor(
		transactable,
		tenantRepository,
		staffService,
		assetService,
	)
	d.AdminAssetInteractor = usecase.NewAdminAssetInteractor(
		assetService,
	)

	d.StaffInteractor = usecase.NewStaffInteractor(
		transactable,
		tenantRepository,
		staffService,
	)

	d.AuthenticationInteractor = usecase.NewAuthenticationInteractor(
		staffAuthenticationRepository,
	)

	d.DebugInteractor = usecase.NewDebugInteractor(
		staffAuthenticationRepository,
	)
}
