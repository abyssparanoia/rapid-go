package create_root_staff_cmd

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
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
		return errors.InternalErr.Wrap(err)
	}
	if email == "" {
		return errors.InternalErr.Errorf("email is required")
	}

	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return errors.InternalErr.Wrap(err)
	}
	if password == "" {
		return errors.InternalErr.Errorf("password is required")
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
