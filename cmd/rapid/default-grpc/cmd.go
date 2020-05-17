package defaultgrpc

import (
	"github.com/spf13/cobra"
)

// NewDefaultGRPCCmd ... new default GRPC cmd
func NewDefaultGRPCCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default-grpc",
		Short: "cli default grpc server",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	cmd.AddCommand(newDefaultGRPCRunCmd())
	return cmd
}

func newDefaultGRPCRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "running default grpc server",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
	return cmd
}
