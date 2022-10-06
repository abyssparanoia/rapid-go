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
	transactable := transactable.NewTransactable()
	authenticationRepository := firebase_repository.NewAuthentication(
		d.FirebaseClient,
	)
	tenantRepository := repository.NewTenant()
	userRepository := repository.NewUser()

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
