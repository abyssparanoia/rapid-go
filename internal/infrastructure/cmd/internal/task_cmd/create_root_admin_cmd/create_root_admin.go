//nolint:forbidigo
package create_root_admin_cmd

import (
	"context"
	"fmt"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/now"
	"github.com/abyssparanoia/rapid-go/internal/pkg/password"
	"github.com/abyssparanoia/rapid-go/internal/usecase"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/spf13/cobra"
)

type CMD struct {
	ctx                 context.Context
	taskAdminInteractor usecase.TaskAdminInteractor
}

func (c *CMD) CreateRootAdmin(cmd *cobra.Command) error {
	email, err := cmd.Flags().GetString("email")
	if err != nil {
		return errors.InternalErr.Wrap(err)
	}
	if email == "" {
		return errors.InternalErr.WithDetail("email is required")
	}

	displayName, err := cmd.Flags().GetString("display-name")
	if err != nil {
		return errors.InternalErr.Wrap(err)
	}
	if displayName == "" {
		return errors.InternalErr.WithDetail("display-name is required")
	}

	// Generate random password
	generatedPassword, err := password.Generate(password.DefaultLength)
	if err != nil {
		return errors.InternalErr.Wrap(err).WithDetail("failed to generate password")
	}

	result, err := c.taskAdminInteractor.Create(
		c.ctx,
		input.NewTaskCreateAdmin(
			email,
			displayName,
			generatedPassword,
			now.Now(),
		),
	)
	if err != nil {
		return err
	}

	// Output result
	fmt.Printf("AdminID: %s\n", result.AdminID)
	fmt.Printf("AuthUID: %s\n", result.AuthUID)
	fmt.Printf("Password: %s\n", result.Password)

	return nil
}
