package http_server_cmd

import (
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/http"
	"github.com/spf13/cobra"
)

func NewHTTPServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http-server",
		Short: "cli http server",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "running http server",
		Run: func(cmd *cobra.Command, args []string) {
			http.Run()
		},
	})
	return cmd
}
