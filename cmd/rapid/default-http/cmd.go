package defaulthttp

import (
	"github.com/spf13/cobra"
)

// NewDefaultHTTPCmd ... new default HTTP cmd
func NewDefaultHTTPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default-http",
		Short: "cli default http server",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	cmd.AddCommand(newDefaultHTTPRunCmd())
	return cmd
}

func newDefaultHTTPRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "running default http server",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
	return cmd
}
