package dependency

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/aws"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito"
	cognito_repository "github.com/abyssparanoia/rapid-go/internal/infrastructure/cognito/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database"
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

	// public
	PublicAuthenticationInteractor usecase.PublicAuthenticationInteractor
	PublicTenantInteractor         usecase.PublicTenantInteractor

	// admin
	AdminTenantInteractor usecase.AdminTenantInteractor
	AdminUserInteractor   usecase.AdminUserInteractor
	AdminAssetInteractor  usecase.AdminAssetInteractor

	// Other
	UserInteractor           usecase.UserInteractor
	AuthenticationInteractor usecase.AuthenticationInteractor
	DebugInteractor          usecase.DebugInteractor
}

func (d *Dependency) Inject(
	ctx context.Context,
	e *environment.Environment,
) {
	d.DatabaseCli = database.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase)

	firebaseCli := firebase.NewClient(e.GCPProjectID)

	gcsCli := gcs.NewClient(ctx)
	gcsBucketHandle := gcs.NewBucketHandle(gcsCli, e.GCPBucketName)

	awsSession := aws.NewSession(e.AWSRegion, e.AWSEmulatorHost)
	_ = s3.NewClient(awsSession)

	cognitoCli := cognito.NewClient(awsSession, e.AWSCognitoEmulatorHost)

	transactable := transactable.NewTransactable()
	_ = firebase_repository.NewAuthentication(
		firebaseCli,
		e.FirebaseClientAPIKey,
	)
	authenticationRepository := cognito_repository.NewAuthentication(
		ctx,
		cognitoCli,
		e.AWSCognitoUserPoolID,
		e.AWSCognitoClientID,
		e.AWSCognitoEmulatorHost,
	)
	tenantRepository := repository.NewTenant()
	userRepository := repository.NewUser()
	assetRepository := gcs_repository.NewAsset(gcsBucketHandle)
	// assetRepository := s3_repository.NewAsset(s3Client, e.AWSBucketName)

	userService := service.NewUser(
		userRepository,
		authenticationRepository,
	)

	d.PublicTenantInteractor = usecase.NewPublicTenantInteractor(
		transactable,
		tenantRepository,
	)

	d.PublicAuthenticationInteractor = usecase.NewPublicAuthenticationInteractor(
		userRepository,
	)

	d.AdminTenantInteractor = usecase.NewAdminTenantInteractor(
		transactable,
		tenantRepository,
	)
	d.AdminUserInteractor = usecase.NewAdminUserInteractor(
		transactable,
		tenantRepository,
		userService,
	)
	d.AdminAssetInteractor = usecase.NewAdminAssetInteractor(
		assetRepository,
	)

	d.UserInteractor = usecase.NewUserInteractor(
		transactable,
		tenantRepository,
		userService,
	)

	d.AuthenticationInteractor = usecase.NewAuthenticationInteractor(
		authenticationRepository,
	)

	d.DebugInteractor = usecase.NewDebugInteractor(
		authenticationRepository,
	)
}
