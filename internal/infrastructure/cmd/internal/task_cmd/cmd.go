package task_cmd

import (
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cmd/internal/task_cmd/create_root_staff_cmd"
	"github.com/spf13/cobra"
)

func NewTaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "cli task",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	cmd.AddCommand(create_root_staff_cmd.NewCreateRootStaffCmd())
	return cmd
}
