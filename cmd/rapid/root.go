package main

import (
	"fmt"
	"os"

	defaulthttp "github.com/abyssparanoia/rapid-go/cmd/rapid/default-http"

	defaultgrpc "github.com/abyssparanoia/rapid-go/cmd/rapid/default-grpc"

	"github.com/spf13/cobra"
)

func newCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rapid",
		Short: "cli tool for rapid-go",
		Long:  "cli tool for rapid-go",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	cmd.AddCommand(defaultgrpc.NewDefaultGRPCCmd())
	cmd.AddCommand(defaulthttp.NewDefaultHTTPCmd())
	return cmd
}

func execute() {
	cmd := newCmdRoot()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
