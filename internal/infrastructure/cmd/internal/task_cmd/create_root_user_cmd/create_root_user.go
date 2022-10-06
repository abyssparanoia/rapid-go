package create_root_user_cmd

import (
	"context"
	"errors"

	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/now"
	"github.com/playground-live/moala-meet-and-greet-back/internal/usecase"
	"github.com/playground-live/moala-meet-and-greet-back/internal/usecase/input"
	"github.com/spf13/cobra"
)

type CMD struct {
	ctx            context.Context
	userInteractor usecase.UserInteractor
}

func (c *CMD) CreateRootUser(cmd *cobra.Command) error {
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

	if err := c.userInteractor.CreateRoot(
		c.ctx,
		input.NewCreateRootUser(
			email,
			password,
			now.Now(),
		),
	); err != nil {
		return err
	}

	return nil
}
