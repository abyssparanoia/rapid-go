package main

import (
	"fmt"
	"os"

	"github.com/abyssparanoia/rapid-go/internal/infrastructure/cmd"
)

func main() {
	cmd := cmd.NewCmdRoot()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
