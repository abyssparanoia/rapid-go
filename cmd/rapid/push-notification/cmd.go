package notification

import (
	"github.com/spf13/cobra"
)

// NewPushNotificationCmd ... new notification server cmd
func NewPushNotificationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push-notification",
		Short: "cli default push notification server",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	cmd.AddCommand(newPushNotificationRunCmd())
	return cmd
}

func newPushNotificationRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "running push notification server",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
	return cmd
}
