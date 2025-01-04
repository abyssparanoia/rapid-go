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
		Use:   "create",
		Short: "create new migration file",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Flags().StringP("name", "n", "", "file name")
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				panic(err)
			}
			if name == "" {
				name = "please_rename_this_file"
			}
			migration.RunNewFile(name)
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "migrate up",
		Run: func(cmd *cobra.Command, args []string) {
			migration.RunUp()
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "extract-schema",
		Short: "extract schema",
		Run: func(cmd *cobra.Command, args []string) {
			migration.RunExtractSchema()
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "sync-constants",
		Short: "sync constants",
		Run: func(cmd *cobra.Command, args []string) {
			migration.RunSyncConstants()
		},
	})
	return cmd
}
