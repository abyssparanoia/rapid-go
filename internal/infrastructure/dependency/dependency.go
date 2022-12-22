package dependency

import (
	"context"

	"firebase.google.com/go/auth"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/transactable"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/firebase"
	firebase_repository "github.com/abyssparanoia/rapid-go/internal/infrastructure/firebase/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/gcs"
	gcs_repository "github.com/abyssparanoia/rapid-go/internal/infrastructure/gcs/repository"
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
	gcsBucketHandle := gcs.NewBucketHandle(gcsCli, e.BucketName)

	transactable := transactable.NewTransactable()
	authenticationRepository := firebase_repository.NewAuthentication(
		d.FirebaseClient,
	)
	tenantRepository := repository.NewTenant()
	userRepository := repository.NewUser()
	assetRepository := gcs_repository.NewAsset(gcsBucketHandle)

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
