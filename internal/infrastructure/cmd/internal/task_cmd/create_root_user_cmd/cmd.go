package create_root_user_cmd

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/dependency"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	"github.com/caarlos0/env/v6"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/spf13/cobra"
)

func NewCreateRootUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-root-user",
		Short: "create root user",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			e := &environment.Environment{}
			if err := env.Parse(e); err != nil {
				panic(err)
			}

			logger := logger.New()
			ctx = ctxzap.ToContext(ctx, logger)

			d := &dependency.Dependency{}
			d.Inject(ctx, e)

			c := &CMD{
				ctx,
				d.UserInteractor,
			}
			if err := c.CreateRootUser(cmd); err != nil {
				panic(err)
			}
		},
	}
	cmd.Flags().StringP("email", "e", "", "email address")
	cmd.Flags().StringP("password", "p", "", "password")
	return cmd
}
