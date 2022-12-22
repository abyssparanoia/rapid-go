package schema_migration_cmd

import (
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cmd/internal/schema_migration_cmd/database_cmd"
	"github.com/spf13/cobra"
)

func NewSchemaMigrationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema-migration",
		Short: "cli schema migration",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	cmd.AddCommand(database_cmd.NewDatabaseCmd())
	return cmd
}
