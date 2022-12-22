package database_cmd

import (
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/database/migration"
	"github.com/spf13/cobra"
)

func NewDatabaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "database",
		Short: "for database",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "migrate up",
		Run: func(cmd *cobra.Command, args []string) {
			migration.RunUp()
		},
	})
	return cmd
}
