package create_root_staff_cmd

import (
	"context"
	"errors"

	"github.com/abyssparanoia/rapid-go/internal/pkg/now"
	"github.com/abyssparanoia/rapid-go/internal/usecase"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/spf13/cobra"
)

type CMD struct {
	ctx             context.Context
	staffInteractor usecase.StaffInteractor
}

func (c *CMD) CreateRootStaff(cmd *cobra.Command) error {
	email, err := cmd.Flags().GetString("email")
	if err != nil {
		return err
	}
	if email == "" {
		return errors.New("email is required")
	}

	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return err
	}
	if password == "" {
		return errors.New("password is required")
	}

	if err := c.staffInteractor.CreateRoot(
		c.ctx,
		input.NewCreateRootStaff(
			email,
			password,
			now.Now(),
		),
	); err != nil {
		return err
	}

	return nil
}
