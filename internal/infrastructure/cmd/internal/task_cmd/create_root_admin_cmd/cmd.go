package create_root_admin_cmd

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/dependency"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	"github.com/caarlos0/env/v11"
	"github.com/spf13/cobra"
)

func NewCreateRootAdminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-root-admin",
		Short: "create root admin",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			e := &environment.Environment{}
			if err := env.Parse(e); err != nil {
				panic(err)
			}

			l := logger.New(e.MinLogLevel)
			ctx = logger.ToContext(ctx, l)

			d := &dependency.Dependency{}
			d.Inject(ctx, e)

			c := &CMD{
				ctx,
				d.TaskAdminInteractor,
			}
			if err := c.CreateRootAdmin(cmd); err != nil {
				panic(err)
			}
		},
	}
	cmd.Flags().StringP("email", "e", "", "email address")
	cmd.Flags().StringP("display-name", "d", "", "display name")
	return cmd
}
