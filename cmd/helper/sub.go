package main

import (
	"github.com/abyssparanoia/rapid-go/cmd/helper/ctxhelper"
	"github.com/spf13/cobra"
)

func newCmdHello() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hello",
		Short: "output hello world",
		Run: func(cmd *cobra.Command, args []string) {
			d := getDeps()
			ctx := ctxhelper.GetContext()
			d.HelperHandler.HelloWorld(ctx)
		},
	}

	return cmd
}
