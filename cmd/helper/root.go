package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// newCmdRoot ...
func newCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "helper",
		Short: "helper cli",
		Long:  `helper cli`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
			}
		}}
	cmd.AddCommand(newCmdHello())
	return cmd
}

func execute() {
	cmd := newCmdRoot()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
