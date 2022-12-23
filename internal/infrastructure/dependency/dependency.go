package dependency

import (
	"context"

	"firebase.google.com/go/auth"
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
	FirebaseClient *auth.Client

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
}

func (d *Dependency) Inject(
	ctx context.Context,
	e *environment.Environment,
) {
	_ = database.NewClient(e.DBHost, e.DBUser, e.DBPassword, e.DBDatabase)

	d.FirebaseClient = firebase.NewClient(e.GCPProjectID)

	gcsCli := gcs.NewClient(ctx)
	gcsBucketHandle := gcs.NewBucketHandle(gcsCli, e.GCPBucketName)

	awsSession := aws.NewSession(e.AWSRegion, e.AWSEmulatorHost)
	_ = s3.NewClient(awsSession)

	cognitoCli := cognito.NewClient(awsSession, e.AWSCognitoEmulatorHost)

	transactable := transactable.NewTransactable()
	_ = firebase_repository.NewAuthentication(
		d.FirebaseClient,
	)
	authenticationRepository := cognito_repository.NewAuthentication(
		ctx,
		cognitoCli,
		e.AWSCognitoUserPoolID,
		e.AWSCognitoEmulatorHost,
	)
	tenantRepository := repository.NewTenant()
	userRepository := repository.NewUser()
	assetRepository := gcs_repository.NewAsset(gcsBucketHandle)
	// assetRepository := s3_repository.NewAsset(s3Client, e.AWSBucketName)

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
		authenticationRepository,
		userRepository,
		tenantRepository,
	)
	d.AdminAssetInteractor = usecase.NewAdminAssetInteractor(
		assetRepository,
	)

	d.UserInteractor = usecase.NewUserInteractor(
		transactable,
		userRepository,
		tenantRepository,
		authenticationRepository,
	)

	d.AuthenticationInteractor = usecase.NewAuthenticationInteractor(
		authenticationRepository,
	)
}
