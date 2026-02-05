package create_root_admin_cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/now"
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

	// Generate random password (16 characters)
	password, err := generatePassword(16)
	if err != nil {
		return errors.InternalErr.Wrap(err).WithDetail("failed to generate password")
	}

	result, err := c.taskAdminInteractor.Create(
		c.ctx,
		input.NewTaskCreateAdmin(
			email,
			displayName,
			password,
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

// generatePassword generates a random password of the specified length
func generatePassword(length int) (string, error) {
	// Calculate bytes needed for base64 encoding
	// base64 encoding expands data by 4/3, so we need length * 3 / 4 bytes
	byteLength := (length * 3) / 4
	if byteLength < 1 {
		byteLength = 1
	}

	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Encode to base64 and truncate to desired length
	encoded := base64.URLEncoding.EncodeToString(bytes)
	if len(encoded) > length {
		encoded = encoded[:length]
	}

	return encoded, nil
}
