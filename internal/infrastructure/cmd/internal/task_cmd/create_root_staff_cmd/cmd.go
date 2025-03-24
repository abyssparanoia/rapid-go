package create_root_staff_cmd

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/dependency"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	"github.com/caarlos0/env/v11"
	"github.com/spf13/cobra"
)

func NewCreateRootStaffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-root-staff",
		Short: "create root staff",
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
				d.StaffInteractor,
			}
			if err := c.CreateRootStaff(cmd); err != nil {
				panic(err)
			}
		},
	}
	cmd.Flags().StringP("email", "e", "", "email address")
	cmd.Flags().StringP("password", "p", "", "password")
	return cmd
}
