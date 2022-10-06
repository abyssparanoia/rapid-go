package cmd

import (
	"github.com/playground-live/moala-meet-and-greet-back/internal/infrastructure/cmd/internal/http_server_cmd"
	"github.com/playground-live/moala-meet-and-greet-back/internal/infrastructure/cmd/internal/task_cmd"
	"github.com/spf13/cobra"
)

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "cli tool for app",
		Long:  "cli tool for app",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	cmd.AddCommand(http_server_cmd.NewHTTPServerCmd())
	cmd.AddCommand(task_cmd.NewTaskCmd())
	return cmd
}
